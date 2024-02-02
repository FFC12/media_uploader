# Media Uploader

This service provides file upload operations for Cloudflare R2 (AWS S3) over WebSocket and utilizes a multi-threaded architecture (more likely concurrent goroutines than threads).

## Getting Started

Follow the steps below to deploy Docker container:

### Run

1. Clone the project:

   ```bash
   git clone https://github.com/FFC12/media_uploader.git
   ```

2. Run the project:
    ```bash
    cd scripts
    chmod u+x docker-deploy.sh
    ./docker-deploy.sh
    ``` 

It'll automatically build Docker image and start the project in Docker. To change configurations for command-line arguments, change it from Dockerfile (.env file can be used in the future for configuration).

## Docker

### Command-Line Arguments

| Argument                     | Default                | Description                                       |
| ---------------------------- | ---------------------- | ------------------------------------------------- |
| `addr`                       | "localhost:8080"       | HTTP service address.                             |
| `perf`                       | false                  | Enable performance testing.                       |
| `workers`                    | 10                     | Number of workers.                                |
| `chBufferSize`               | 100                    | Channel buffer size.                              |
| `workerMemoryLimit`          | 75                     | Worker memory limit in percentage.                              |
| `useSpawnerWithMemoryLimit`  | true                   | Use worker spawner with memory limit.             |
| `enableSimpleInterface`      | false                  | Enable simple interface to upload files.         |
| `saveUploadsTemporarily`     | false                  | Save uploaded files temporarily.                 |

### Usage Example

Here is an example of how to run the project with custom configurations:

```bash
./media_uploader_binary -addr="0.0.0.0:8080" -perf=true -workers=20 -chBufferSize=200 -workerMemoryLimit=100 -useSpawnerWithMemoryLimit=false -enableSimpleInterface=true -saveUploadsTemporarily=true
```

## Using Simple Interface

The argument `enableSimpleInterface` is used to enable the simple interface. The application offers two main file upload methods:

### 1. Stream Upload

To perform a stream upload (it will use the webcam), follow these steps:

- Navigate to the stream upload endpoint to test: [http://localhost:8080/stream](http://localhost:8080/stream) 
- Start streaming your file data to the WebSocket.

### 2. File Select Upload

For file select upload (it will ask for a file), follow these steps:

- Go to the file select upload endpoint to test: [http://localhost:8080/file_select](http://localhost:8080/file_select) 
- Choose a file using the provided interface and initiate the upload.

