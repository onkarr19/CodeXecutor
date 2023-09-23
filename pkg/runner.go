package pkg

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// LanguageConfig defines the image name and command for each supported language.
type LanguageConfig struct {
	ImageName string
	Cmd       []string
}

// ExecuteCodeInContainer runs the provided code in a Docker container with the specified language dependency.
func executeCodeInContainer(code string, languageConfig LanguageConfig) (string, string, error) {

	// Initialize Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", "", err
	}

	// Define host configuration
	hostConfig := &container.HostConfig{
		Resources: container.Resources{
			Memory: int64(250 * 1024 * 1024), // 250 MB memory limit
		},
		// Disable internet access
		NetworkMode: "none",
	}

	// Set container configuration options
	containerConfig := &container.Config{
		Image: languageConfig.ImageName,
		Cmd:   languageConfig.Cmd,
	}

	// Wait for the container to finish, with a timeout of 5 seconds
	timeout := time.Second * 5
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create a Docker container
	fmt.Println("Creating Docker container...")
	resp, err := cli.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig,
		nil,
		nil,
		"",
	)
	if err != nil {
		return "", resp.ID, err
	}

	// Start the Docker container
	fmt.Println("Starting Docker container...")
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", resp.ID, err
	}

	waitCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case <-waitCh:
	case err := <-errCh:
		if err != nil {
			return "", resp.ID, err
		}
	case <-ctx.Done():
		return "", resp.ID, fmt.Errorf("container execution timed out after %d seconds", timeout/time.Second)
	}

	// Retrieve container logs
	fmt.Println("Retrieving container logs...")
	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return "", resp.ID, err
	}
	defer out.Close()

	// Read and format the container output
	containerOutput, err := io.ReadAll(out)
	if err != nil {
		return "", resp.ID, err
	}
	return strings.TrimSpace(string(containerOutput)), resp.ID, nil
}

// remove the container
func removeDockerContainer(containerID string) error {
	fmt.Println("Removing Docker container...")
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	if err := cli.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{Force: true}); err != nil {
		log.Printf("Error removing Docker container: %v\n", err)
		return err
	}
	return nil
}

func ExecuteAndCleanupContainer(code, language string) (string, error) {

	languageConfigs := map[string]LanguageConfig{
		"C": {
			ImageName: "c_image_name_here",
			Cmd:       []string{"gcc", "your_c_program.c", "-o", "output", "&&", "./output"},
		},
		"C++": {
			ImageName: "gcc:latest",
			Cmd:       []string{"sh", "-c", "echo \"$CODE\" > code.cpp && g++ code.cpp -o output && ./output", code},
		},
		"Python": {
			ImageName: "python:3.9",
			Cmd:       []string{"python", "-c", code},
		},
		"Java": {
			ImageName: "java_image_name_here",
			Cmd:       []string{"java", "-jar", "your_jar_file.jar"},
		},
		"JavaScript": {
			ImageName: "nodejs_image_name_here",
			Cmd:       []string{"node", "your_javascript_program.js"},
		},
		"Go": {
			ImageName: "golang_image_name_here",
			Cmd:       []string{"go", "run", "your_golang_program.go"},
		},
	}

	languageConfig, ok := languageConfigs[language]
	if !ok {
		return "", fmt.Errorf("language %s is not supported", language)
	}

	// execute the container
	output, containerID, err := executeCodeInContainer(code, languageConfig)
	// remove the container
	removeDockerContainer(containerID)
	if err != nil {
		return "", err
	}

	return output, nil
}

func main() {
	// Example usage
	code :=
		`print("Hello World!")
`
	language := "Python"

	output, err := ExecuteAndCleanupContainer(code, language)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Output: %s\n", output)
	}
}
