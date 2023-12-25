package redis

import (
	"CodeXecutor/models"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr            string
	Password        string
	DB              int
	PoolSize        int
	MinIdleConns    int
	MaxRetries      int
	MinRetryBackoff int
	MaxRetryBackoff int
}

// LoadRedisConfig loads the Redis configuration from a TOML file
func LoadRedisConfig(filePath string) (*RedisConfig, error) {
	var cfg RedisConfig
	if _, err := toml.DecodeFile(filePath, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// ConnectRedis creates a Redis connection pool based on the provided configuration
func ConnectRedis() *redis.Client {

	cfg, err := LoadRedisConfig("config/redis.toml")
	if err != nil {
		panic(fmt.Sprintf("Error loading Redis config: %s", err.Error()))
	}

	// Create a new Redis Options struct
	options := &redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	// Create a pool of Redis connections
	poolSize := 10
	minIdleConns := 5

	pool := redis.NewClient(&redis.Options{
		Addr:         options.Addr,
		Password:     options.Password,
		DB:           options.DB,
		PoolSize:     poolSize,
		MinIdleConns: minIdleConns,
	})

	// Check if the connection to Redis is successful
	_, err = pool.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %s", err.Error()))
	} else {
		fmt.Println("Connected to Redis")
	}

	return pool
}

func EnqueueItem(client *redis.Client, queueName string, codeSubmission models.Job) error {
	// Convert the codeSubmission struct to a JSON string
	result, err := json.Marshal(codeSubmission)
	if err != nil {
		return err
	}

	// Enqueue the JSON string in the Redis list
	return client.LPush(context.Background(), queueName, result).Err()
}

func DequeueItem(client *redis.Client, queueName string) (models.Job, error) {
	// Dequeue the JSON string from the Redis list
	result, err := client.BRPop(context.Background(), 0, queueName).Result()
	if err != nil {
		return models.Job{}, err
	}

	value := result[1]

	// Convert the JSON string back to the codeSubmission struct
	var codeSubmission models.Job
	err = json.Unmarshal([]byte(value), &codeSubmission)
	if err != nil {
		return models.Job{}, err
	}

	return codeSubmission, nil
}

func SetCache(client *redis.Client, key string, data interface{}, expiration time.Duration) error {
	// Convert the data to a JSON string
	result, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Set the JSON string in Redis with expiration time
	return client.Set(context.Background(), key, result, expiration).Err()
}
