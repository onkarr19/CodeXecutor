package redis

import (
	"CodeXecutor/models"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
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

type Config struct {
	Redis RedisConfig `toml:"redis"`
}

var (
	ConfigSingle *Config
	once         sync.Once
	pool         *redis.Client
	err          error
)

// LoadRedisConfig loads the Redis configuration from a TOML file
func LoadRedisConfig(filePath string) (*Config, error) {
	var config Config

	// Read the TOML file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return &config, err
	}

	// Unmarshal the TOML data into the Config struct
	err = toml.Unmarshal(data, &config)
	if err != nil {
		return &config, err
	}

	return &config, nil
}

// ConnectRedis creates a Redis connection pool based on the provided configuration
func ConnectRedis() *redis.Client {

	once.Do(func() {
		ConfigSingle, err = LoadRedisConfig("config/redis.toml")
		if err != nil {
			panic(fmt.Sprintf("Error loading Redis config: %s", err.Error()))
		}

		// Create a new Redis Options struct
		options := &redis.Options{
			Addr:     ConfigSingle.Redis.Addr,
			Password: ConfigSingle.Redis.Password,
			DB:       ConfigSingle.Redis.DB,
		}

		// Create a pool of Redis connections
		poolSize := 10
		minIdleConns := 5

		pool = redis.NewClient(&redis.Options{
			Addr:         options.Addr,
			Password:     options.Password,
			DB:           options.DB,
			PoolSize:     poolSize,
			MinIdleConns: minIdleConns,
		})

		// Check if the connection to Redis is successful
		_, err = pool.Ping(context.Background()).Result()
		if err != nil {
			panic(fmt.Sprintf("Failed to connect to Redis: %s\n", err.Error()))
		} else {
			fmt.Println("Connected to Redis")
		}
	})

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

func GetCache(client *redis.Client, key string) (string, error) {
	// Check if the key exists in the cache
	result, err := client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		// Key does not exist in the cache
		return "", fmt.Errorf("key not found in cache")
	} else if err != nil {
		// Error occurred while fetching from the cache
		return "", err
	}

	// Key found in the cache, return the result
	return result, nil
}
