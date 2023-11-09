package redis

import (
	"CodeXecutor/models"
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {
	// Initialize and return a connection to your Redis instance
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password by default
		DB:       0,                // Default redis database
	})

	// Check if the connection to Redis is successful
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}

	return client
}

func EnqueueCodeSubmission(client *redis.Client, queueName string, codeSubmission models.Job) error {
	// Convert the codeSubmission struct to a JSON string
	codeSubmissionJSON, err := json.Marshal(codeSubmission)
	if err != nil {
		return err
	}

	// Enqueue the JSON string in the Redis list
	return client.LPush(context.Background(), queueName, codeSubmissionJSON).Err()
}
