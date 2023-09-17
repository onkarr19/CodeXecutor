package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	fmt.Println("in here")
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("Worker %d started job %d\n", id, job)
		time.Sleep(10 * time.Second)

		// Create a unique container name
		containerName := fmt.Sprintf("mycontainer%d", job)

		// Create a context with a timeout
		ctx, cancel := context.WithCancel(context.Background())

		// Start a goroutine to handle the timeout and cancel the context if necessary
		go func() {
			time.Sleep(5 * time.Second) // Adjust the timeout duration as needed
			cancel()
		}()

		// Build and run the Docker command with the context
		cmd := exec.CommandContext(ctx, "docker", "run", "--name", containerName, "python:3.8", "python", "-c", "import time; time.sleep(10)")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				fmt.Printf("Job %d timed out after 5 seconds\n", job)
			} else {
				fmt.Printf("Error running Docker container: %v\n", err)
			}
		} else {
			// Remove the container after it has run
			removeCmd := exec.Command("docker", "rm", "-f", containerName)
			if err := removeCmd.Run(); err != nil {
				fmt.Printf("Error removing Docker container: %v\n", err)
			}
		}

		fmt.Printf("Worker %d finished job %d\n", id, job)
	}
}

func main() {
	// /* Sample Execution of workers

	const numWorkers = 5
	const numJobs = 13

	jobs := make(chan int, numJobs)
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		fmt.Println("sending ", i)
		go worker(i, jobs, &wg)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}

	// Close the jobs channel to signal that no more jobs will be added
	close(jobs)

	// Wait for all workers to finish
	wg.Wait()

	// */

}
