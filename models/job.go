package models

type Job struct {
	ID       string `json:"id"`       // unique identifier
	Language string `json:"language"` // Programming language used in the code
	Code     string `json:"code"`     // The user's code
	Time     int    `json:"time"`     // The time of submission
}

// DockerConfig represents the configuration for the Docker container.
type DockerConfig struct {
	Image    string `json:"image"`    // Docker image name, e.g., "python:3.9"
	Code     string `json:"code"`     // User's code to be executed
	Language string `json:"language"` // Programming language (used for selecting the image)
}
