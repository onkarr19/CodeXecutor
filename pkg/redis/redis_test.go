package redis

import (
	"CodeXecutor/models"
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testConfigFilePath = "test-redis-config.toml"

func TestLoadRedisConfig(t *testing.T) {
	// Create a sample Redis configuration file for testing
	testConfigContent := `
		[redis]
		Addr = "localhost:6379"
		Password = ""
		DB = 0
		PoolSize = 10
		MinIdleConns = 5
	`

	err := os.WriteFile(testConfigFilePath, []byte(testConfigContent), 0644)
	assert.NoError(t, err, "Error creating test Redis config file")
	defer os.Remove(testConfigFilePath)

	config, err := LoadRedisConfig(testConfigFilePath)
	assert.NoError(t, err, "Error loading Redis config")
	assert.NotNil(t, config, "Config should not be nil")
	assert.Equal(t, "localhost:6379", config.Redis.Addr, "Incorrect Redis address in config")

}

func TestConnectRedis(t *testing.T) {
	// Ensure that ConnectRedis initializes the Redis client and connects to Redis
	client := ConnectRedis()
	assert.NotNil(t, client, "Redis client should not be nil")

	// Check if the client can perform a basic operation (Ping)
	pong, err := client.Ping(context.Background()).Result()
	assert.NoError(t, err, "Failed to ping Redis")
	assert.Equal(t, "PONG", pong, "Unexpected ping response")
}

func TestEnqueueDequeueItem(t *testing.T) {
	client := ConnectRedis()

	// Create a sample job for testing
	job := models.Job{
		ID:       "123",
		Language: "go",
		Code:     "fmt.Println('Hello, World!')",
		Time:     int(time.Now().Unix()),
	}

	// Enqueue the job
	err := EnqueueItem(client, "test-queue", job)
	assert.NoError(t, err, "Error enqueuing item")

	// Dequeue the job
	dequeuedJob, err := DequeueItem(client, "test-queue")
	assert.NoError(t, err, "Error dequeuing item")
	assert.Equal(t, job, dequeuedJob, "Enqueued and dequeued jobs should be equal")

}

func TestSetGetCache(t *testing.T) {
	client := ConnectRedis()

	// Set data in cache
	key := "unique-key-2"
	data := models.CompilationResult{ExitCode: 0, Output: ".", Error: nil}
	err := SetCache(client, key, data, time.Minute)
	assert.NoError(t, err, "Error setting data in cache")

	// Get data from cache
	cachedData, err := GetCache(client, key)
	assert.NoError(t, err, "Error getting data from cache")
	assert.NotEmpty(t, cachedData, "Cached data should not be empty")

	// compare the cached data
	assert.NoError(t, err, "Error unmarshalling cached data")
	assert.Equal(t, data, cachedData, "Cached data does not match expected data")
}

func TestSetCacheExpiredGetCache(t *testing.T) {
	client := ConnectRedis()
	defer client.Close()

	key := "expiredKey"
	data := "testData"
	expiration := time.Millisecond * 100 // Very short expiration time for testing

	// Test SetCache with short expiration time
	err := SetCache(client, key, data, expiration)
	assert.NoError(t, err, "SetCache should not return an error")

	// Wait for expiration
	time.Sleep(expiration + time.Millisecond*50)

	// Test GetCache for an expired key
	_, err = GetCache(client, key)
	assert.Error(t, err, "GetCache should return an error for an expired key")
}

func TestGetCacheKeyNotFound(t *testing.T) {
	client := ConnectRedis()
	defer client.Close()

	// Test GetCache for a non-existent key
	_, err := GetCache(client, "nonExistentKey")
	assert.Error(t, err, "GetCache should return an error for a non-existent key")
}
