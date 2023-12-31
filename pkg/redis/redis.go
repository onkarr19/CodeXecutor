package redis

import (
	"CodeXecutor/models"
	"CodeXecutor/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	clientPool   *redis.Client
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
		configPath, er := utils.GetFilePath("config", "redis.toml")
		if er != nil {
			panic(er)
		}

		ConfigSingle, err = LoadRedisConfig(configPath)
		if err != nil {
			log.Fatalf("Error loading Redis config: %v", err)
		}

		// Create a clientPool of Redis connections
		poolSize := 10
		minIdleConns := 5

		// Create a new Redis Options struct
		options := &redis.Options{
			Addr:         ConfigSingle.Redis.Addr,
			Password:     ConfigSingle.Redis.Password,
			DB:           ConfigSingle.Redis.DB,
			PoolSize:     poolSize,
			MinIdleConns: minIdleConns,
		}

		clientPool = redis.NewClient(options)

		// Check if the connection to Redis is successful
		_, err = clientPool.Ping(context.Background()).Result()
		if err != nil {
			log.Fatalf("Failed to connect to Redis: %v\n", err)
		} else {
			log.Println("Connected to Redis")
		}
	})

	return clientPool
}

func EnqueueItem(queueName string, codeSubmission models.Job) error {
	// Convert the codeSubmission struct to a JSON string
	result, err := json.Marshal(codeSubmission)
	if err != nil {
		return err
	}

	// Enqueue the JSON string in the Redis list
	return clientPool.LPush(context.Background(), queueName, result).Err()
}

func DequeueItem(queueName string) (models.Job, error) {
	// Dequeue the JSON string from the Redis list
	result, err := clientPool.BRPop(context.Background(), 0, queueName).Result()
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

func SetCache(key string, data interface{}, expiration time.Duration) error {
	// Convert the data to a JSON string
	result, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	// Set the JSON string in Redis with expiration time
	return clientPool.Set(context.Background(), key, result, expiration).Err()
}

func GetCache(key string) (models.CompilationResult, error) {
	// Check if the key exists in the cache
	result, err := clientPool.Get(context.Background(), key).Result()
	if err == redis.Nil {
		// Key does not exist in the cache
		return models.CompilationResult{}, fmt.Errorf("key not found in cache")
	} else if err != nil {
		// Error occurred while fetching from the cache
		return models.CompilationResult{}, err
	}

	var compilationResult models.CompilationResult
	err = json.Unmarshal([]byte(result), &compilationResult)
	if err != nil {
		return models.CompilationResult{}, err
	}

	// Key found in the cache, return the result
	return compilationResult, nil
}
