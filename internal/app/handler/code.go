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

type Response struct {
	Message      string `json:"message"`
	SubmissionId string `json:"submissionid"`
}

// HandleResponse writes a JSON response to the http.ResponseWriter.
func HandleResponse(w http.ResponseWriter, status int, err error, result models.CompilationResult) {

	response := map[string]interface{}{
		"found": err != redis.Nil,
		"data":  nil,
	}

	if err != nil && err != redis.Nil {
		// Handle Redis error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if response["found"].(bool) {
		// Key found
		response["data"] = map[string]interface{}{
			"output":   result.Output,
			"error":    result.Error,
			"exitcode": result.ExitCode,
		}
	}

	sendJSONResponse(w, response, http.StatusOK)
}

// HandleError handles HTTP errors and writes a JSON response with an error message.
func HandleError(w http.ResponseWriter, status int, err error) {
	log.Printf("Error: %v", err)
	HandleResponse(w, status, err, models.CompilationResult{})
}

// HandleSubmissionResponse handles the response after submitting code.
func HandleSubmissionResponse(w http.ResponseWriter, key string, status int) {
	message := models.CompilationResult{}

	// Wait for 0.5 second (500 milliseconds)
	time.Sleep(500 * time.Millisecond)

	// Retrieve result if available
	result, err := redisClient.GetCache(key)
	if err == nil {
		// Set the response code and message based on the result
		status = http.StatusOK
		message = result
	}

	// final response
	HandleResponse(w, status, nil, message)
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
	err = redisClient.EnqueueItem(queueName, job)
	if err != nil {
		log.Printf("Failed to enqueue code submission: %v", err)
		// Handle the error and respond to the user with an error message
		http.Error(w, "Failed to submit code", http.StatusInternalServerError)
		return
	}

	HandleSubmissionResponse(w, job.ID, http.StatusAccepted)
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

// HandleResult handles requests to retrieve code processing results.
func HandleResult(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	result, err := redisClient.GetCache(key)
	if err != nil {
		HandleError(w, http.StatusNoContent, err)
		return
	}

	HandleResponse(w, http.StatusOK, err, result)
}

func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.WriteHeader(statusCode)

	// Encode and write the response
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
