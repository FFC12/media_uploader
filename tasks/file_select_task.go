package tasks

import (
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/media_uploader/core"
)

// FileSelectUpload represents a task for streaming file uploads.
type FileSelectUpload struct {
	task core.Task
	Conn *websocket.Conn
}

// Execute method implements the task execution logic for selected file uploads.
func (t *FileSelectUpload) Execute() error {
	// Close the connection when the task execution is complete.
	defer t.Conn.Close()

	// Create a 'temp' folder if it does not exist.
	_, err := os.Stat("temp")
	if os.IsNotExist(err) {
		errDir := os.Mkdir("temp", 0755)
		if errDir != nil {
			fmt.Println("Error:", errDir, "not critical. It keeps going.")
		}
	}

	// Generate a unique filename using UUID.
	uniqueFileName := uuid.New().String()

	// Create a binary file to store the uploaded data.
	binaryFile, err := os.Create("temp/" + uniqueFileName + ".bin")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer binaryFile.Close()

	// Read and write data in chunks until "EOF" is received.
	for {
		_, message, err := t.Conn.ReadMessage()
		if err != nil {
			// Handle normal closure, check file size, and cleanup if necessary.
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				break
			}

			// Check file size
			stat, err := os.Stat("temp/" + uniqueFileName + ".bin")
			if err != nil {
				fmt.Println(err)
				return err
			}

			// Check file size <1kb, remove incomplete file if true.
			if stat.Size() < 1024 {
				os.Remove("temp/" + uniqueFileName + ".bin")
				return nil
			}

			return errors.New("socket has been closed - this is not a critical error")
		}

		// Break the loop when "EOF" is received.
		if string(message) == "EOF" {
			break
		}

		// Write the received data to the binary file.
		_, err = binaryFile.Write(message)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	// Read the file format from the connection.
	_, format, err := t.Conn.ReadMessage()
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Rename the binary file to include the file format in the filename.
	newFileName := "temp/" + uniqueFileName + "." + string(format)
	err = os.Rename("temp/"+uniqueFileName+".bin", newFileName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Send a WebSocket close message.
	err = t.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
