package tests

import (
	"CodeXecutor/models"
	"CodeXecutor/pkg/redis"
	"CodeXecutor/utils"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestFullFlow tests the full end-to-end flow of enqueueing, dequeuing, and caching.
func TestFullFlow(t *testing.T) {
	client := redis.ConnectRedis()

	// Create a sample job for testing
	job := models.Job{
		ID:       "456",
		Language: "python",
		Code:     "print('Hello, Redis!')",
		Time:     int(time.Now().Unix()),
	}

	// Enqueue the job
	err := redis.EnqueueItem(client, "test-queue", job)
	assert.NoError(t, err, "Error enqueuing item")

	// Dequeue the job
	dequeuedJob, err := redis.DequeueItem(client, "test-queue")
	assert.NoError(t, err, "Error dequeuing item")
	assert.Equal(t, job, dequeuedJob, "Enqueued and dequeued jobs should be equal")

	// Set and get data in cache
	key := utils.GenerateUniqueID()
	data := map[string]interface{}{"status": "success"}
	err = redis.SetCache(client, key, data, time.Minute)
	assert.NoError(t, err, "Error setting data in cache")

	cachedData, err := redis.GetCache(client, key)
	assert.NoError(t, err, "Error getting data from cache")
	assert.NotEmpty(t, cachedData, "Cached data should not be empty")

	// Unmarshal the cached data and compare
	var result map[string]interface{}
	err = json.Unmarshal([]byte(cachedData), &result)
	assert.NoError(t, err, "Error unmarshalling cached data")
	assert.Equal(t, data, result, "Cached data does not match expected data")
}
