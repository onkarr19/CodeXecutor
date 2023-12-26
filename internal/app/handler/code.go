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

	response := struct {
		SubmissionId string `json:"submissionid"`
		Message      string `json:"message"`
	}{}
	response.SubmissionId = job.ID

	// Wait for 0.5 second (500 milliseconds)
	time.Sleep(500 * time.Millisecond)

	result, err := redisClient.GetCache(client, job.ID)
	if err != nil {
		// Respond to the user with a message indicating that the code has been submitted
		response.Message = "Your code has been submitted for processing."
		w.WriteHeader(http.StatusAccepted) // 202 Accepted status code

	} else {
		// return result
		response.Message = result
		w.WriteHeader(http.StatusOK) // 200 Ok status code
	}

	// Convert the response structure to JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write the JSON response to the response writer
	w.Write(responseJSON)
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
		http.Error(w, "No data found", http.StatusNoContent)
		return
	}
	response.Message = result

	// Convert the response structure to JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the HTTP status code to 200
	w.WriteHeader(http.StatusOK)

	// Write the JSON response to the response writer
	w.Write(responseJSON)
}
