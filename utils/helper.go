package utils

import (
	"errors"
	"path/filepath"
	"runtime"

	"github.com/google/uuid"
)

// GenerateUniqueID generates a unique ID using Google's UUID package
func GenerateUniqueID() string {
	return uuid.New().String()
}

// GetFilePath returns the path of the file
func GetFilePath(dir, file string) (string, error) {
	// Get the absolute path of the directory containing the source file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("failed to get the caller's filename")
	}

	// Construct the absolute path to the file
	filePath := filepath.Join(filepath.Dir(filename), "..", dir, file)

	return filePath, nil
}
