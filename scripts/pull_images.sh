#!/bin/bash

# Define the list of Docker images to pull
images=("gcc:10.3" "python:3.9" "openjdk:11.0.12" "node:14.17" "golang:1.21")

# Loop through the list and pull each image
for image in "${images[@]}"; do
  echo "Pulling image: $image"
  docker pull -q "$image"
done

echo "All images have been pulled successfully!"
