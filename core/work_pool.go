package core

import (
	"fmt"
	"sync"
)

// WorkPool represents a simple worker pool implementation.
type WorkPool struct {
	wg        sync.WaitGroup
	workerNum int
	taskCh    chan Task
}

// NewPool creates a new WorkPool with the specified number of workers and buffer size for tasks.
func NewPool(workerNumber int, bufferSize int) *WorkPool {
	return &WorkPool{
		workerNum: workerNumber,
		taskCh:    make(chan Task, bufferSize),
	}
}

// Done decrements the internal WaitGroup counter.
func (wp *WorkPool) Done() {
	wp.wg.Done()
}

// Wait waits for all goroutines to finish.
func (wp *WorkPool) Wait() {
	wp.wg.Wait()
}

// Run submits a task to the worker pool.
func (wp *WorkPool) Run(task Task) {
	wp.wg.Add(1)
	go func(t Task) {
		defer wp.Done()

		// Perform cleanup and exit if an error occurs
		defer func() {
			if r := recover(); r != nil {
				// Handle the panic or cleanup as needed
				wp.handleError(fmt.Errorf("panic: %v", r))
			}
		}()

		// Execute the task and handle errors
		if err := t.Execute(); err != nil {
			wp.handleError(err)
			// Perform any additional cleanup here if needed
		}
	}(task)
}

// Start initializes and starts the worker pool.
func (wp *WorkPool) Start() {
	wp.wg.Add(wp.workerNum) // Increment WaitGroup for each worker
	for i := 0; i < wp.workerNum; i++ {
		go func(id int) {
			// Print worker id when starting
			fmt.Printf("Worker %d started\n", id)
			defer wp.wg.Done()
			for task := range wp.taskCh {
				task.Execute()
			}
		}(i)
	}
}

// handleError is a simple function to handle errors.
func (wp *WorkPool) handleError(err error) {
	// Handle the error here
	fmt.Println("Error:", err)
}

// Close closes the task channel and waits for all workers to finish.
func (wp *WorkPool) Close() {
	close(wp.taskCh)
	wp.Wait() // Wait for all workers to finish before returning
}
