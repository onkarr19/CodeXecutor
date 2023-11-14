package worker

import (
	"context"
	"log"

	"CodeXecutor/models"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// GenerateDockerContainer dynamically generates a Docker container for code execution.
func GenerateDockerContainer(ctx context.Context, client *client.Client, config models.DockerConfig) (string, error) {
	// Create a container configuration
	containerConfig := &container.Config{
		Image:        config.Image,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		OpenStdin:    true,
		StdinOnce:    true,
		Cmd:          []string{config.Language, "-c", config.Code}, // Customize this based on the image and code execution method
	}

	// Create a container host configuration
	hostConfig := &container.HostConfig{}

	// Create the container
	resp, err := client.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "my-coding-platform-container")
	if err != nil {
		log.Printf("Error creating container: %v\n", err)
		return "", err
	}

	// Start the container
	if err := client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Printf("Error starting container: %v\n", err)
		return "", err
	}

	return resp.ID, nil
}

// StopAndRemoveContainer stops and removes a Docker container.
func StopAndRemoveContainer(ctx context.Context, client *client.Client, containerID string) error {
	timeout := int(10) // Adjust the timeout as needed

	stopOptions := container.StopOptions{
		Timeout: &timeout,
	}

	// Stop the container
	if err := client.ContainerStop(ctx, containerID, stopOptions); err != nil {
		log.Printf("Error stopping container: %v\n", err)
		return err
	}

	// Remove the container
	if err := client.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{}); err != nil {
		log.Printf("Error removing container: %v\n", err)
		return err
	}

	return nil
}
