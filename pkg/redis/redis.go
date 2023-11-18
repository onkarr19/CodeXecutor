package redis

import (
	"CodeXecutor/config"
	"CodeXecutor/models"
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {
	cfg := config.NewRedisConfig()
	// Initialize and return a connection to your Redis instance
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
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

func DequeueCodeSubmission(client *redis.Client, queueName string) (models.Job, error) {
	// Dequeue the JSON string from the Redis list
	codeSubmissionJSON, err := client.RPop(context.Background(), queueName).Result()
	if err != nil {
		return models.Job{}, err
	}

	// Convert the JSON string back to the codeSubmission struct
	var codeSubmission models.Job
	err = json.Unmarshal([]byte(codeSubmissionJSON), &codeSubmission)
	if err != nil {
		return models.Job{}, err
	}

	return codeSubmission, nil
}
