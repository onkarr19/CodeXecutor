# Golang Dockerfile
FROM golang:1.21.5

# Install Docker
RUN apt-get update && apt-get install -y docker.io

# Set the working directory inside the container
WORKDIR /app

# Copy the Go project files into the container
COPY . .

# Build the Go application
RUN go build -o main ./cmd

# Expose the port your application is running on
EXPOSE 8080

# Command to run your application
CMD ["./main"]