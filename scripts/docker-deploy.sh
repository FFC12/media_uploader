#!/usr/bin/bash

cd ../

# Build the Docker image
docker build --pull --rm -f "Dockerfile" -t mediauploader:latest "."

# Run the Docker container in detached mode
docker run -d -p 8080:8080 -i -t mediauploader
