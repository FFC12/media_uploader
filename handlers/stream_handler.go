package handlers

import (
	"fmt"
	"net/http"

	tasks "github.com/media_uploader/tasks"
)

// Http stream handler
func StreamHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	task := &tasks.StreamUploadTask{
		Conn:                   conn,
		SaveUploadsTemporarily: SaveUploadsTemporarily,
	}

	WorkerPool.Run(task)

	// fligramTask := &tasks.FligramStamp{
	// 	Image: "Somethin which is not fligram",
	// }

	// WorkerPool.Run(fligramTask)
}
