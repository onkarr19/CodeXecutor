package config

// RedisConfig holds the configuration parameters for connecting to a Redis server.
type RedisConfig struct {
	Addr     string // Redis server address (e.g., "localhost:6379")
	Password string // Redis password (leave empty for no password)
	DB       int    // Redis database (0 by default)
}

// NewRedisConfig returns a default RedisConfig with some default values.
func NewRedisConfig() *RedisConfig {
	return &RedisConfig{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
}

// RedisQueueName is the name of the Redis queue for code submissions.
const RedisQueueName = "code-submissions"
