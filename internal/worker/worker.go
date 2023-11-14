package worker

import (
	"CodeXecutor/models"
	"context"
	"log"
	"sync"
)

// // Job represents a code compilation job.
// type Job struct {
// 	// Add job-related fields here, e.g., code, language, etc.
// }

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
		case _, ok := <-w.jobQueue:
			if !ok {
				// Job queue has been closed, exit the worker
				return
			}

			// Handle the code compilation job here
			containerID, err := GenerateDockerContainer(w.ctx, nil, models.DockerConfig{})
			if err != nil {
				log.Printf("Error generating Docker container: %v\n", err)
				// Handle the error appropriately
				continue
			}

			// Wait for code execution to complete (you may implement this logic)
			// ...

			// Remove the Docker container
			if err := StopAndRemoveContainer(w.ctx, nil, containerID); err != nil {
				log.Printf("Error stopping and removing Docker container: %v\n", err)
				// Handle the error appropriately
				continue
			}

		case <-w.ctx.Done():
			{
				// Context canceled, exit the worker
				return
			}
		}
	}
}
