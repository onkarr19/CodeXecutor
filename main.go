package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

const redisHost = "localhost:6379"
const redisPassword = ""

type Submission struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	Lang     string `json:"lang"`
	Expected string `json:"expected"`
}

type Verdict struct {
	ID      string `json:"id"`
	Output  string `json:"output"`
	Verdict bool   `json:"verdict"`
}

func RedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       0, // Use the default Redis database
	})
}

func getStatus(redisClient *redis.Client, key string) int {
	// Get the integer value associated with the key from Redis
	val, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return -1
	}

	// Convert the retrieved string to an integer
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return -1
	}

	return intVal
}

func publishToQueue(redisClient *redis.Client, queueName string, key string, data interface{}) error {

	// Serialize the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Push the serialized JSON to the queue
	err = redisClient.LPush(ctx, queueName, jsonData).Err()
	if err != nil {
		return err
	}

	// Set a key-value pair
	err = redisClient.Set(ctx, key, 1, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func popFromInputQueue(redisClient *redis.Client, queueName string) (Submission, error) {
	// Pop the serialized JSON from the right end of the queue
	jsonData, err := redisClient.RPop(ctx, queueName).Result()
	if err != nil {
		return Submission{}, err
	}

	// Deserialize the JSON into a Submission struct
	var data Submission
	err = json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return Submission{}, err
	}

	return data, nil
}

func serializeVerdict(verdict Verdict) (string, error) {
	verdictJSON, err := json.Marshal(verdict)
	if err != nil {
		return "", err
	}
	return string(verdictJSON), nil
}

func deserializeVerdict(verdictJSON string) (Verdict, error) {
	var v Verdict
	err := json.Unmarshal([]byte(verdictJSON), &v)
	return v, err
}

func setVerdictInRedis(redisClient *redis.Client, key string, verdict Verdict) error {
	verdictJSON, err := serializeVerdict(verdict)
	if err != nil {
		return err
	}
	return redisClient.HSet(ctx, "verdicts", key, verdictJSON).Err()
}

func getVerdictFromRedis(redisClient *redis.Client, key string) (Verdict, error) {
	verdictJSONString, err := redisClient.HGet(ctx, "verdicts", key).Result()
	if err != nil {
		return Verdict{}, err
	}
	return deserializeVerdict(verdictJSONString)
}

func worker(id int, wg *sync.WaitGroup, redisClient *redis.Client) {
	defer wg.Done()
	for {
		// Pop a message from the Redis List
		submission, err := popFromInputQueue(redisClient, "inputqueue")
		if err != nil {
			if err == redis.Nil {
				// No message available, worker will wait for new messages
				// fmt.Printf("Worker %d: No message available, waiting...\n", id)
				// break // Exit the loop and stop the worker
			} else {
				// Handle other Redis-related errors here, e.g., reconnect to Redis or log the error
				fmt.Printf("Worker %d: Error: %v\n", id, err)
				continue
			}
		} else {

			// Now you can access the 'ID' field of 'submission'
			// fmt.Printf("Worker %d: Received message: %s\n", id, submission.ID)
			err := redisClient.Incr(ctx, submission.ID).Err()
			if err != nil {
				fmt.Println("Error incrementing status:", err)
			}

			time.Sleep(3 * time.Second)
			// TODO: Write docker logic here.
			output := "456"
			verdict := Verdict{ID: submission.ID, Output: output}

			if submission.Expected != "" {
				verdict.Verdict = output == submission.Expected
			}

			err = setVerdictInRedis(redisClient, verdict.ID, verdict)
			if err != nil {
				fmt.Println("Error setting verdict:", err)
				return
			}

			err = redisClient.Incr(ctx, verdict.ID).Err()
			if err != nil {
				fmt.Println("Error incrementing status:", err)
			}

			// fmt.Printf("Worker %d: Completed: %s\n", id, submission)
		}
	}
}

func submitHandler(w http.ResponseWriter, r *http.Request, redisClient *redis.Client) {
	// generate unique ID
	// publishToQueue input queue, set status=queued
	// key := "abcd1234"
	// value := Submission{ID: key, Code: "this is C++ code", Lang: "C++"}
	// publishToQueue(redisClient, "inputqueue", key, value)

	fmt.Fprintln(w, "Job received and queued.")
}

func statusHandler(w http.ResponseWriter, r *http.Request, redisClient *redis.Client) {
	// Get the key from the request, for example, by reading a query parameter
	id := r.URL.Query().Get("id")

	// Check if the key is empty
	if id == "" {
		http.Error(w, "Key parameter is required", http.StatusBadRequest)
		return
	}

	// Get the status from Redis using the getStatus function
	status := getStatus(redisClient, id)

	// Return the status as a response
	fmt.Fprintf(w, "Status for key %s: %d", id, status)
}

func resultHandler(w http.ResponseWriter, r *http.Request, redisClient *redis.Client) {

	id := r.URL.Query().Get("id")

	// Set the Content-Type header to indicate that the response is in JSON format
	w.Header().Set("Content-Type", "application/json")

	// Get the Verdict from Redis
	verdict, err := getVerdictFromRedis(redisClient, id)

	if err != nil {
		http.Error(w, "Error getting verdict from Redis", http.StatusInternalServerError)
		return
	}

	// Marshal the verdict struct to JSON
	verdictJSON, err := json.Marshal(verdict)
	if err != nil {
		http.Error(w, "Error marshaling verdict to JSON", http.StatusInternalServerError)
		return
	}

	// Write the JSON response to the HTTP response writer
	w.WriteHeader(http.StatusOK)
	w.Write(verdictJSON)
}

func main() {
	/* Sample Execution of workers

	const numWorkers = 5
	const numJobs = 13

	jobs := make(chan int, numJobs)
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		fmt.Println("sending ", i)
		go worker(i, jobs, &wg)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}

	// Close the jobs channel to signal that no more jobs will be added
	close(jobs)

	// Wait for all workers to finish
	wg.Wait()

	*/

	// Create a redis client
	redisClient := RedisClient()

	// // Worker pool configuration
	const numWorkers = 5
	// jobs := make(chan int, numWorkers)
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, &wg, redisClient)
	}

	// Create a Gorilla Mux router
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		submitHandler(w, r, redisClient)
	}).Methods("POST")

	router.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		statusHandler(w, r, redisClient)
	}).Methods("GET")

	router.HandleFunc("/result", func(w http.ResponseWriter, r *http.Request) {
		resultHandler(w, r, redisClient)
	}).Methods("GET")

	// // Start the HTTP server
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8080", nil))

	// Wait for all workers to finish.
	wg.Wait()

	// Close the Redis client when done
	// if err := redisClient.Close(); err != nil {
	// 	fmt.Printf("Error closing Redis client: %v\n", err)
	// }
}
