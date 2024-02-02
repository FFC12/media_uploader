package handlers

import (
	"github.com/gorilla/websocket"
	wp "github.com/media_uploader/core"
)

// upgrader is a WebSocket upgrader with specified read and write buffer sizes.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 256,
}

// WorkerPool is a global instance of the worker pool used by handlers.
// It is initialized with 10 workers and a buffer size of 100 for task channel.
var WorkerPool *wp.WorkerPool

// WorkerSpawner is a global instance of the worker spawner with 75% memory limit based on runtime.
// This gives us control for spawning workers based on memory usage as much as possible.
// It's better for memory usage based workers than using a fixed number of workers (in some cases).
var WorkerSpawner *wp.WorkerSpawnerWithMemoryLimit

var SaveUploadsTemporarily = false

func InitializeWorkerConfig(workerCount, chBufferSize int, memoryLimit uint64) {
	WorkerPool = wp.NewPool(workerCount, chBufferSize)
	WorkerSpawner = wp.NewWorkerSpawnerWithMemoryLimit(memoryLimit)
}
