package handler

import (
	"CodeXecutor/models"
	redisClient "CodeXecutor/pkg/redis"
	"CodeXecutor/utils"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func init() {
	client = redisClient.ConnectRedis()
}

// HandleCodeSubmission handles incoming code submissions.
func HandleCodeSubmission(w http.ResponseWriter, r *http.Request) {
	// Extract code submission data from the request
	job, err := extractCodeSubmission(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Specify the Redis queue name for code submissions
	queueName := "code-submissions"

	// Enqueue the code submission in Redis for processing
	err = redisClient.EnqueueItem(client, queueName, job)
	if err != nil {
		log.Printf("Failed to enqueue code submission: %v", err)
		// Handle the error and respond to the user with an error message
		http.Error(w, "Failed to submit code", http.StatusInternalServerError)
		return
	}

	// Wait for 1 second (1000 milliseconds)
	time.Sleep(500 * time.Millisecond)

	result, err := redisClient.GetCache(client, job.ID)
	responseMessage := ""
	if err != nil {
		// Respond to the user with a message indicating that the code has been submitted
		responseMessage = "Your code has been submitted for processing. Submission ID: " + job.ID
	} else {
		// return result
		responseMessage = result
	}
	w.WriteHeader(http.StatusAccepted) // 202 Accepted status code
	w.Write([]byte(responseMessage))
}

func extractCodeSubmission(r *http.Request) (models.Job, error) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		// Handle the error and return an appropriate response
		return models.Job{}, err
	}

	var job models.Job
	if err := json.Unmarshal(body, &job); err != nil {
		log.Printf("Failed to unmarshal JSON: %v", err)
		// Handle the error and return an appropriate response
		return models.Job{}, err
	}

	job.ID = utils.GenerateUniqueID()
	return job, nil
}

func HandleResult(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	response := struct {
		Message string `json:"message"`
	}{}

	result, err := redisClient.GetCache(client, key)
	if err != nil {
		log.Println(err)
		http.Error(w, "No data found", http.StatusNotFound)
		return
	}
	response.Message = result

	// Convert the response structure to JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to indicate JSON content
	w.Header().Set("Content-Type", "application/json")

	// Set the HTTP status code to 412 (Precondition Failed)
	w.WriteHeader(http.StatusPreconditionFailed)

	// Write the JSON response to the response writer
	w.Write(responseJSON)
}
