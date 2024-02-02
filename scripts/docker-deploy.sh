#!/usr/bin/bash
cd ../
docker build --pull --rm -f "Dockerfile" -t mediauploader:latest "."
docker run -d -p 8080:8080 -i -t mediauploader
docker exec -i mediauploader bash -c "cd app && ./server -useSpawnerWithMemoryLimit=true -enableSimpleInterface=true -addr=0.0.0.0:8080"