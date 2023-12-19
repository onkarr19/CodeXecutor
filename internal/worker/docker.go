package worker

import (
	"fmt"
	"log"
	"time"

	"CodeXecutor/models"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

// GenerateAndStartContainer dynamically generates a Docker container for code execution.
func (w *Worker) GenerateAndStartContainer(config models.DockerConfig) (string, error) {
	containerConfig := &container.Config{
		Image:        config.Image,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		OpenStdin:    true,
		StdinOnce:    true,
		Cmd:          []string{config.Language, "-c", config.Code},
	}

	hostConfig := &container.HostConfig{}

	resp, err := w.client.ContainerCreate(w.ctx, containerConfig, hostConfig, nil, nil, time.Now().Format("20060102150405"))
	if err != nil {
		log.Printf("Error creating container: %v\n", err)
		return "", err
	}

	if err := w.client.ContainerStart(w.ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Printf("Error starting container: %v\n", err)
		return "", err
	}

	// Wait for the container to finish
	waitResultCh, errCh := w.client.ContainerWait(w.ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case waitResult := <-waitResultCh:
		if waitResult.StatusCode != 0 {
			log.Printf("Container exited with non-zero status code: %d\n", waitResult.StatusCode)
			return resp.ID, fmt.Errorf("container exited with non-zero status code: %d", waitResult.StatusCode)
		}
	case err := <-errCh:
		if err != nil {
			log.Printf("Error waiting for container to finish: %v\n", err)

			return resp.ID, err
		}
	case <-w.ctx.Done():
		log.Printf("Context canceled while waiting for container to finish.\n")
		return resp.ID, w.ctx.Err()
	}

	// // Capture container output
	// out, err := w.client.ContainerLogs(w.ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	// if err != nil {
	// 	log.Printf("Error retrieving container logs: %v\n", err)
	// 	return resp.ID, "", err
	// }
	// defer out.Close()

	// // Read the output
	// output, err := io.ReadAll(out)
	// if err != nil {
	// 	log.Printf("Error reading container logs: %v\n", err)
	// 	return "", "", err
	// }

	// log.Printf("Container output:\n%s\n", output)

	return resp.ID, nil
}

// StopAndRemoveContainer stops and removes a Docker container.
func (w *Worker) StopAndRemoveContainer(containerID string) error {
	timeout := int(10)

	stopOptions := container.StopOptions{
		Timeout: &timeout,
	}

	if err := w.client.ContainerStop(w.ctx, containerID, stopOptions); err != nil {
		log.Printf("Error stopping container: %v\n", err)
		return err
	}

	if err := w.client.ContainerRemove(w.ctx, containerID, types.ContainerRemoveOptions{}); err != nil {
		log.Printf("Error removing container: %v\n", err)
		return err
	}

	return nil
}
