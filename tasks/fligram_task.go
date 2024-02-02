package tasks

import (
	"os"
	"time"

	"github.com/media_uploader/core"
)

// FileSelectUpload represents a task for streaming file uploads.
type FligramStamp struct {
	task  core.Task
	Image string
}

// Execute method implements the task execution logic for selected file uploads.
func (t *FligramStamp) Execute() error {
	// Fake task execution for 5 seconds.
	time.Sleep(5 * time.Second)
	os.WriteFile("stamp.txt", []byte(t.Image), 0644)
	return nil
}
