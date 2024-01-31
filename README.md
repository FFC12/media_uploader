# Media Uploader

This service provides file upload operations over WebSocket and utilizes a multi-threaded architecture (more likely concurrent goroutines than threads).

## Getting Started

Follow the steps below to start the project. To set worker counts, edit the `workerCount` variable in the `handlers/common.go` file.
  
### Run

1. Clone the project:

   ```bash
   git clone https://github.com/FFC12/media_uploader.git
   ```

2. Run the project:
    ```bash
    cd media_uploader
    ```

3. Run the application:
    ```bash
    go run main.go
    ```

## Usage

The application offers two main file upload methods:

### 1. Stream Upload

To perform a stream upload (it will use the webcam), follow these steps:

- Navigate to the stream upload endpoint to test: [http://localhost:8080/stream](http://localhost:8080/stream)
- Connect to the WebSocket endpoint: [ws://localhost:8080/upload_stream](ws://localhost:8080/upload_stream)
- Start streaming your file data to the WebSocket.

### 2. File Select Upload

For file select upload (it will ask for a file), follow these steps:

- Go to the file select upload endpoint to test: [http://localhost:8080/file_select](http://localhost:8080/file_select)
- Connect to the WebSocket endpoint: [ws://localhost:8080/upload_select_file](ws://localhost:8080/upload_select_file)
- Choose a file using the provided interface and initiate the upload.

Note that the WebSocket endpoints should be used for the corresponding upload methods. Make sure to handle the upload process according to the selected method.
