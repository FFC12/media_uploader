package handlers

import (
	"fmt"
	"net/http"

	tasks "github.com/media_uploader/tasks"
)

// Http selected file handler
func FileSelectHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	task := &tasks.FileSelectUpload{
		Conn: conn,
	}

	WorkPool.Run(task)
}
