package handler

import (
	"CodeXecutor/config"
	"CodeXecutor/models"
	"CodeXecutor/pkg/redis"
	"CodeXecutor/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// HandleCodeSubmission handles incoming code submissions.
func HandleCodeSubmission(w http.ResponseWriter, r *http.Request) {
	// Extract code submission data from the request
	job, err := extractCodeSubmission(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(job)

	// Enqueue the code submission in Redis for processing
	client := redis.ConnectRedis()
	defer client.Close()

	// Specify the Redis queue name for code submissions
	queueName := config.RedisQueueName

	err = redis.EnqueueItem(client, queueName, job)
	if err != nil {
		log.Printf("Failed to enqueue code submission: %v", err)
		// Handle the error and respond to the user with an error message
		http.Error(w, "Failed to submit code", http.StatusInternalServerError)
		return
	}

	// log.Printf("Code submission with ID %s enqueued for processing", job.ID)

	// Respond to the user with a message indicating that the code has been submitted
	responseMessage := "Your code has been submitted for processing. Submission ID: " + job.ID
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
		fmt.Println(string(body))
		return models.Job{}, err
	}

	job.ID = utils.GenerateUniqueID()
	return job, nil
}
