package main

import (
	"CodeXecutor/internal/app"
	"CodeXecutor/internal/worker"
	"context"
	"log"
)

func main() {
	ctx := context.Background()

	// Initialize the application server
	server := app.NewServer()
	server.Start()

	// Initialize the worker pool with min and max worker limits
	workerPool := worker.NewWorkerPool(ctx, 1, 4)
	defer workerPool.Stop()

	// Wait for termination signal
	<-ctx.Done()

	// Graceful shutdown
	server.Stop()

	log.Println("Server gracefully stopped")
}
