# CodeXecutor - Code Executor with Docker and Redis

This project provides an HTTP server in Go that enables you to run code in an isolated Docker environment with specified constraints. The server communicates with Redis to manage code submissions, execution, and result retrieval. The Docker containers have no access to the internet, a time limit of 5 seconds, and are automatically removed after execution.


## Project Structure
```
CodeXecutor/
├── README.md       # Project documentation and instructions
├── cmd/            # Command-line application code
│   └── main.go     # Main application entry point
├── config/
│   └── redis.go    # Redis Configuration 
├── go.mod          # Go module file (dependency management)
├── go.sum          # Go dependencies checksum file
├── internal/       # Internal application code
│   ├── app/        # Application-specific code
│   │   ├── handler/
│   │   │   └── code.go  # Code handling logic
│   │   └── server.go    # Application server code
├── models/         # Data models
│   ├── job.go      # Job-related data models
├── pkg/            # Reusable packages and libraries
│   ├── redis/      # Redis-related code
│   │   └── redis.go     # Redis connection code
├── tmp/            # Temporary files or logs
│   ├── build-errors.log  # Build error logs
│   └── main/        # Compiled application or binary
└── utils/          # Utility code
    └── helper.go   # Helper functions and utilities
```