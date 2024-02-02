package tasks

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gorilla/websocket"
	uploader "github.com/media_uploader/amazon"
	"github.com/media_uploader/core"
)

// StreamUploadTask represents a task for streaming file uploads.
type StreamUploadTask struct {
	task                   core.Task
	Conn                   *websocket.Conn
	SaveUploadsTemporarily bool

	// Add mutex to protect shared resources
	mu sync.Mutex
}

// Execute method implements the task execution logic for streaming file uploads.
func (t *StreamUploadTask) Execute() error {
	// Close the connection when the task execution is complete.
	defer t.Conn.Close()

	// Create a 'temp' folder if it does not exist.
	if t.SaveUploadsTemporarily {
		_, err := os.Stat("temp")
		if os.IsNotExist(err) {
			errDir := os.Mkdir("temp", 0755)
			if errDir != nil {
				fmt.Println("Error:", errDir, "not critical. It keeps going.")
			}
		}
	}

	// Read first chunk for video data
	_, data, err := t.Conn.ReadMessage()
	if err != nil {
		core.LogError("Error (while reading first chunk)", err)
		return err
	}

	// Deserialize first chunk
	firstChunk, err := JsonSerializer.Deserialize(data)
	if err != nil {
		core.LogError("Error (while deserializing first chunk)", err)
		return err
	}

	// Extract the MIME type from the first chunk
	mimeType := strings.Split(firstChunk.MimeType, "/")[1]

	// Extansion for video
	extension := strings.Split(mimeType, ";")[0]

	// Generate a unique filename using UUID.
	uniqueFileName := firstChunk.MediaId
	fileName := uniqueFileName + "." + extension

	// Create a binary file to store the uploaded data.
	var binaryFile *os.File

	// Create a binary file to store the uploaded data.
	if t.SaveUploadsTemporarily {
		binaryFile, err = os.Create("temp/" + fileName)
		if err != nil {
			core.LogError("Error (while creating binary file)", err)
			return err
		}
		defer binaryFile.Close()
	}

	context := context.Background()

	t.mu.Lock()

	var completedParts = make([]types.CompletedPart, 0)
	var eTags = make([]string, 0)
	var partNumber int32 = 1
	var loc string = ""
	var svc *s3.Client
	var resp *s3.CreateMultipartUploadOutput

	var buffer []byte
	var directUploadFlag bool = true
	var multipartUploadFlag bool = false

	// Read and write data in chunks until "EOF" is received.
	for {
		_, message, err := t.Conn.ReadMessage()
		if err != nil {
			// Handle normal closure, check file size, and cleanup if necessary.
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				break
			}

			if t.SaveUploadsTemporarily {
				// Check file size
				stat, err := os.Stat("temp/" + fileName)
				if err != nil {
					core.LogError("Error (while checking file size)", err)
					return err
				}

				// Check file size <1kb, remove incomplete file if true.
				if stat.Size() < 1024 {
					os.Remove("temp/" + fileName)
					return nil
				}
			}

			return errors.New("socket has been closed - sync failed")
		}

		// Break the loop when "EOF" is received.
		if string(message) == "EOF" {
			break
		}

		// Write the received data to the binary file.
		if t.SaveUploadsTemporarily {
			_, err = binaryFile.Write(message)
			if err != nil {
				core.LogError("Error (while writing to binary file)", err)
				return err
			}
		}

		buffer = append(buffer, message...)

		// NOTE: Amazon S3 mandates a minimum part size of 5 MB for multipart uploads.
		// Our approach is to upload in 5 MB parts if the buffer size exceeds this threshold.
		// Otherwise, we upload the data in a single part.
		// However, this strategy imposes a 50 GB upper limit (5 * 10000 MiB) on the data size
		// due to the minimum part size.
		// In a scenario where 200 users upload <5 MB data (e.g., 3 MB each) and 800 users upload >5 MB data (e.g., 250 MB each),
		// the calculated RAM usage for 1000 connections is as follows:
		// (3 * 100) MB + (800 * 5) MB = ~4.3 GB.
		// It's important to note that Goroutines have different lifetimes, and during processing,
		// some reserved memory will be released by the Go garbage collector,
		// especially when handling smaller uploads like the <5MB MB example.
		if len(buffer) >= (1024*1024)*5 {
			if !multipartUploadFlag {
				svc, resp, err = uploader.StreamUploadInit(&context, mimeType, fileName)
				if err != nil {
					core.LogError("Error (while initializing multipart upload)", err)
					return err
				}
				multipartUploadFlag = true
				directUploadFlag = false
			}

			uploadResult, err := uploader.StreamUpload(&context, svc, resp, buffer, partNumber)
			if err != nil {
				core.LogError("Error (while uploading part): %s", err)
				return err
			}

			var numb int32 = int32(partNumber)
			completedParts = append(completedParts, types.CompletedPart{
				ETag:       uploadResult.ETag,
				PartNumber: &numb,
			})

			partNumber += 1
			buffer = []byte{}
		}
	}

	if !directUploadFlag && multipartUploadFlag {
		// Check completed parts tag and part number
		if len(eTags) > 0 {
			for i := 0; i < len(eTags); i++ {
				core.LogInfo("======  ETags ======")
				core.LogInfo(fmt.Sprintf("ETag: %s", eTags[i]))
				core.LogInfo(fmt.Sprintf("Part number: %d", completedParts[i].PartNumber))
				core.LogInfo("====== !ETags ======")
			}
		}

		loc, err = uploader.StreamDone(&context, svc, resp, completedParts)
		if err != nil {
			fmt.Println(err)
			return err
		}
		core.LogInfo(fmt.Sprintf("Video uploaded successfully. Location: %s", loc))
	} else if directUploadFlag {
		loc, err = uploader.DirectUpload(&context, mimeType, fileName, buffer)
		if err != nil {
			fmt.Println(err)
			return err
		}
		core.LogInfo(fmt.Sprintf("Video uploaded successfully. Location: %s", loc))
	} else {
		core.LogError("Error (while uploading video): %s", errors.New("failed to upload video"))
	}

	t.mu.Unlock()

	err = t.Conn.WriteMessage(websocket.TextMessage, []byte(loc))
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Send a WebSocket close message.
	err = t.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Upload completed"))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
