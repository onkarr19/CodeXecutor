# DockerWiz - Code Runner with Docker and Redis

This project provides an HTTP server in Go that enables you to run code in an isolated Docker environment with specified constraints. The server communicates with Redis to manage code submissions, execution, and result retrieval. The Docker containers have no access to the internet, a time limit of 5 seconds, and are automatically removed after execution.

## Features

- Workers subscribe to the Redis list, execute submitted code, and store the result in Redis.
- Submit the code code using `/submit` `POST` endpoint.
- Check the status of code execution using the `/status?id=1` endpoint.
- Retrieve the execution result using the `/result?id=1` endpoint.

<!-- ## Usage -->

### Installation

1. Make sure you have Docker installed on your machine.
2. Clone this repository.

### Getting Started

1. Navigate to the `cmd` directory:

```bash
cd DockerWiz/cmd
```

2. Open the `main.go` file and replace the `code` and `language` variables with your desired values.

3. Run the application:

```bash
go run main.go
```

This will start the HTTP server, along with workers.

## Project Structure

```
DockerWiz/
├── README.md
├── go.mod
├── go.sum
├── cmd/
│   └── main.go
├── pkg/
│   └── runner.go
└── vendor/
```

- `cmd`: Contains the main application entry point.
- `pkg`: Contains the code for running code in Docker.
<!-- - `vendor`: (Optional) Contains vendor dependencies. -->

## Contributing

If you'd like to contribute to this project, please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Make your changes and commit them with descriptive messages.
4. Push your branch to your fork.
5. Create a pull request.