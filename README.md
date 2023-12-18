# CodeXecutor - Code Executor with Docker and Redis

This project provides an HTTP server in Go that enables you to run code in an isolated Docker environment with specified constraints. The server communicates with Redis to manage code submissions, execution, and result retrieval. The Docker containers have no access to the internet, a time limit of 5 seconds, and are automatically removed after execution.


## Project Structure
```
CodeXecutor/
├── Makefile                    # Project dependencies
├── README.md                   # Project documentation and instructions
├── cmd/                        # Command-line application code
│   └── main.go                 # Main application entry point
├── config/
│   └── redis.go                # Redis Configuration 
├── go.mod                      # Go module file (dependency management)
├── go.sum                      # Go dependencies checksum file
├── internal/                   # Internal application code
│   ├── app/                    # Application-specific code
│   │   ├── handler/
│   │   │   └── code.go         # Code handling logic
│   │   └── server.go           # Application-server code
│   └── worker/
│       ├── docker.go           # Docker container logic
│       ├── worker.go           # Worker-specific code
│       └── workerpool.go       # Workerpool management
├── models/
│   ├── job.go                  # Job-related data models
├── pkg/
│   ├── redis/                  # Redis-related code
│   │   └── redis.go            # Redis connection code
├── tmp/
│   ├── build-errors.log        # Build error logs
│   └── main/                   # Compiled application or binary
└── utils/
    └── helper.go               # Helper functions and utilities
```

## Usage
### Starting Dependencies
```bash
make start-services
```
This command initiates the necessary services for the project to function properly. Make sure you have Docker and Docker Compose installed on your system.

### Running the Server


```bash
run run ./cmd
```
This will start the server component of the project.


### Stopping Dependencies
```bash
make stop-services
```

Executing this command will gracefully shut down the services, ensuring a proper termination of the project environment.
