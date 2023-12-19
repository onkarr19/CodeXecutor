package worker

import (
	"CodeXecutor/models"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Worker represents a worker that handles code compilation jobs.
type Worker struct {
	ctx      context.Context
	jobQueue <-chan models.Job
	client   *client.Client
	// Add other worker-related fields here
}

// NewWorker creates a new Worker instance.
func NewWorker(jobQueue <-chan models.Job) *Worker {
	ctx := context.Background()
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil
	}
	return &Worker{ctx: ctx, jobQueue: jobQueue, client: dockerClient}
}

// Start starts the worker to handle jobs.
func (w *Worker) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case job, ok := <-w.jobQueue:
			if !ok {
				// Job queue has been closed, exit the worker
				return
			}

			w.handleJob(job)

		case <-w.ctx.Done():
			{
				// Context canceled, exit the worker
				return
			}
		}
	}
}

func (w *Worker) handleJob(job models.Job) {
	// Map programming languages to corresponding Docker images
	languageToImage := map[string]string{
		"cpp":    "gcc:10.3",
		"python": "python:3.9",
		"java":   "openjdk:11.0.12",
		"node":   "node:14.17",
		"golang": "golang:1.21",
		// Add more languages and their corresponding images as needed
	}

	// Check if the provided language is supported
	image, ok := languageToImage[job.Language]
	if !ok {
		log.Printf("Unsupported programming language: %s\n", job.Language)
		// Handle the error appropriately
		return
	}

	containerID, err := w.GenerateAndStartContainer(models.DockerConfig{
		Image:    image,
		Language: job.Language,
		Code:     job.Code,
	})

	if err != nil {
		log.Printf("Error generating Docker container: %v\n", err)
		// Handle the error appropriately
	}

	// Retrieve container logs
	logs, err := w.getContainerLogs(containerID)
	if err != nil {
		log.Printf("Error getting container logs: %v\n", err)
		// Handle the error appropriately
	} else {
		fmt.Println("Logs:", logs)
	}

	// Remove the Docker container
	if err := w.StopAndRemoveContainer(containerID); err != nil {
		log.Printf("Error stopping and removing Docker container: %v\n", err)
		// Handle the error appropriately
		return
	}
}

func (w *Worker) getContainerLogs(containerID string) (string, error) {
	var logsBuffer bytes.Buffer
	out, err := w.client.ContainerLogs(w.ctx, containerID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return "", err
	}
	defer out.Close()
	io.Copy(&logsBuffer, out)
	return logsBuffer.String(), nil
}
