package utils

import "github.com/google/uuid"

// GenerateUniqueID generates a unique ID using Google's UUID package
func GenerateUniqueID() string {
	return uuid.New().String()
}
