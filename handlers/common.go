package handlers

import (
	"github.com/gorilla/websocket"
	wp "github.com/media_uploader/core"
)

// upgrader is a WebSocket upgrader with specified read and write buffer sizes.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// workerCount is the number of workers in the worker pool.
var workerCount int = 3

// chBufferSize is the buffer size for the task channel.
var chBufferSize int = 100

// WorkPool is a global instance of the worker pool used by handlers.
// It is initialized with 10 workers and a buffer size of 100 for task channel.
var WorkPool = wp.NewPool(workerCount, chBufferSize)
