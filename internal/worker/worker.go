package worker

import (
	"CodeXecutor/models"
	"context"
	"log"
	"sync"
)

// Worker represents a worker that handles code compilation jobs.
type Worker struct {
	ctx      context.Context
	jobQueue <-chan models.Job
	// Add other worker-related fields here
}

// NewWorker creates a new Worker instance.
func NewWorker(ctx context.Context, jobQueue <-chan models.Job) *Worker {
	return &Worker{
		ctx:      ctx,
		jobQueue: jobQueue,
		// Initialize other fields and dependencies
	}
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
	// Handle the code compilation job here
	containerID, err := GenerateAndStartContainer(w.ctx, nil, models.DockerConfig{
		Image:    "your_docker_image", // Provide the actual Docker image name
		Language: job.Language,
		Code:     job.Code,
	})

	if err != nil {
		log.Printf("Error generating Docker container: %v\n", err)
		// Handle the error appropriately
		return
	}

	// Remove the Docker container
	if err := StopAndRemoveContainer(w.ctx, nil, containerID); err != nil {
		log.Printf("Error stopping and removing Docker container: %v\n", err)
		// Handle the error appropriately
		return
	}
}
