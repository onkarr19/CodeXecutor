package worker

import (
	"CodeXecutor/models"
	"CodeXecutor/pkg/redis"
	"context"
	"log"
	"sync"
	"time"
)

// WorkerPool represents a dynamic pool of workers.
type WorkerPool struct {
	minWorkers int
	maxWorkers int
	jobQueue   chan models.Job
	workers    []*Worker
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	// Add other worker pool-related fields and dependencies here
}

// NewWorkerPool initializes and returns a new WorkerPool instance.
func NewWorkerPool(ctx context.Context, minWorkers, maxWorkers int) *WorkerPool {
	jobQueue := make(chan models.Job)
	ctx, cancel := context.WithCancel(ctx)

	wp := &WorkerPool{
		minWorkers: minWorkers,
		maxWorkers: maxWorkers,
		jobQueue:   jobQueue,
		ctx:        ctx,
		cancel:     cancel,
		// Initialize other fields and dependencies
	}

	wp.initWorkers()

	return wp
}

// initWorkers starts the minimum number of workers.
func (wp *WorkerPool) initWorkers() {
	wp.startWorkers(wp.minWorkers)
}

// AddWorkers adds new workers to the pool.
func (wp *WorkerPool) AddWorkers(count int) {
	wp.startWorkers(count)
}

// startWorkers starts the specified number of workers.
func (wp *WorkerPool) startWorkers(count int) {
	workersToAdd := min(count, wp.maxWorkers-len(wp.workers))

	for i := 0; i < workersToAdd; i++ {
		w := NewWorker(wp.jobQueue)
		wp.workers = append(wp.workers, w)
		wp.wg.Add(1)
		go w.Start(&wp.wg)
	}
}

// SubmitJob submits a job to the worker pool.
func (wp *WorkerPool) SubmitJob(job models.Job) {
	wp.jobQueue <- job
}

func PullData(wp *WorkerPool, queueName string) {
	client := redis.ConnectRedis()
	defer client.Close()
	for {
		// Dequeue item from Redis queue
		job, err := redis.DequeueItem(client, queueName)
		if err != nil {
			log.Println("Error dequeueing item from Redis:", err)
			continue
		}

		// Submit the job to the worker pool
		wp.SubmitJob(job)
	}
}

// Stop stops the worker pool and all workers.
func (wp *WorkerPool) Stop() {
	close(wp.jobQueue)
	wp.wg.Wait()
	wp.cancel()
	log.Println("Worker pool stopped")
}

// MonitorSystemLoad periodically monitors system load and adjusts worker count.
func (wp *WorkerPool) MonitorSystemLoad() {
	for {
		select {
		case <-wp.ctx.Done():
			return
		default:
			// Implement system load monitoring logic here
			// You can access metrics like CPU usage or queue length to make decisions

			// Example: Check if the queue is too long, and add more workers if needed
			if len(wp.jobQueue) > 10 && len(wp.workers) < wp.maxWorkers {
				wp.AddWorkers(1)
			}

			// Example: Check if the queue is empty, and remove workers if there are more than the minimum
			if len(wp.jobQueue) == 0 && len(wp.workers) > wp.minWorkers {
				wp.RemoveWorkers(1)
			}

			// Sleep for some time before checking again
			time.Sleep(5 * time.Second)
		}
	}
}

// RemoveWorkers removes workers from the pool.
func (wp *WorkerPool) RemoveWorkers(count int) {
	for i := 0; i < count; i++ {
		if len(wp.workers) > 0 {
			// w := wp.workers[len(wp.workers)-1]
			wp.workers = wp.workers[:len(wp.workers)-1]
			// w.Stop(&wp.wg) // Implement a Stop method in Worker to gracefully stop them
		}
	}
}
