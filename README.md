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
│   ├── middleware
│   │   └── json.go             # set content-type to json
│   └── worker/
│       ├── docker.go           # Docker container logic
│       ├── worker.go           # Worker-specific code
│       └── workerpool.go       # Workerpool management
├── models/
│   ├── job.go                  # Job-related data models
│   └── output.go               # output data model
├── pkg/
│   ├── redis/                  # Redis-related code
│   │   └── redis.go            # Redis connection code
├── scripts/
│   ├── pull_images.sh          # required docker images 
├── tmp/
│   ├── build-errors.log        # Build error logs
│   └── main/                   # Compiled application or binary
└── utils/
    └── helper.go               # Helper functions and utilities
```

## Usage

### Download docker images
Make sure to give execute permissions to the script:
```bash
chmod +x scripts/pull_images.sh
./scripts/pull_images.sh
```


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


### API endpoints
##### Submit Code
```http
POST /submit
```
Endpoint for submitting code.

Request:
Body: 
```json
{
    "Code": "import time;print(time.time());time.sleep(1);print(264/0)",
    "language": "python",
    "time": {{currentTimestamp}} // optional
}
```
Example
```bash
curl -X POST -H "Content-Type: application/json" -d '{
    "Code": "import time;print(time.time());time.sleep(1);print(264/0)",
    "language": "python",
    "time": 1640588800
}' http://localhost:8080/submit
```


##### Get Result:
```bash
GET /result?key={submissionKey}
```
Endpoint for retrieving the result of a submitted code.


Request:
Parameters:
key (string, required): The unique key associated with the code submission.

Example:
```bash
curl http://localhost:8080/result?key=d3389ac4-1080-47c9-b326-19d8437afc2a
```

Response:
```json
{
    "message": "{\"ExitCode\":0,\"Output\":\"1703569908.9141312\\r\\n26.4\\r\\n\",\"Error\":null}"
}
```

Example:
```json
{
    "message": "{\"ExitCode\":0,\"Output\":\"1703569908.9141312\\r\\n26.4\\r\\n\",\"Error\":null}"
}
```

### Stopping Dependencies
```bash
make stop-services
```

Executing this command will gracefully shut down the services, ensuring a proper termination of the project environment.
