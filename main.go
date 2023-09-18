package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

const redisHost = "localhost:6379"
const redisPassword = ""

type Submission struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Lang string `json:"lang"`
}

type TestOutput struct {
	ID     string `json:"id"`
	Output string `json:"output"`
}

func RedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       0, // Use the default Redis database
	})
}

func publishToQueue(redisclient *redis.Client, queueName string, key string, data interface{}) error {

	// Serialize the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Push the serialized JSON to the queue
	err = redisclient.LPush(ctx, queueName, jsonData).Err()
	if err != nil {
		return err
	}

	return nil
}

func worker(id int, wg *sync.WaitGroup, redisClient *redis.Client) {
	defer wg.Done()
	for {
		// Pop a message from the Redis List
		message, err := redisClient.RPop(ctx, "inputqueue").Result()
		if err != nil {
			if err == redis.Nil {
				// No message available, worker will wait for new messages
				fmt.Printf("Worker %d: No message available, waiting...\n", id)
				// break // Exit the loop and stop the worker
			} else {
				// Handle other Redis-related errors here, e.g., reconnect to Redis or log the error
				fmt.Printf("Worker %d: Error: %v\n", id, err)
				continue
			}
		}

		if message != "" {
			fmt.Printf("Worker %d: Received message: %s\n", id, message)
			time.Sleep(3 * time.Second)
			// TODO: Write docker logic here.

			output := TestOutput{ID: "123", Output: "456"}
			publishToQueue(redisClient, "outputqueue", "key", output)
			fmt.Printf("Worker %d: Completed: %s\n", id, message)
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

func statusHandler(w http.ResponseWriter, r *http.Request) {
	// return status of key in redis input queue
	fmt.Fprintln(w, "Job received and queued.")
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	// return result of code execution

	fmt.Fprintln(w, "Job received and queued.")
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

	// // Define routes
	router.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		submitHandler(w, r, redisClient)
	}).Methods("POST")
	router.HandleFunc("/status", statusHandler).Methods("GET")
	router.HandleFunc("/result", resultHandler).Methods("GET")

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
