package core

import (
	"fmt"
	"runtime"
	"sync"
)

// WorkerSpawnerWithMemoryLimit represents a worker pool with dynamic memory-based limits.
type WorkerSpawnerWithMemoryLimit struct {
	wg           sync.WaitGroup
	taskCh       chan Task
	stopCh       chan struct{}
	maxMemoryPct uint64 // Maximum allowed memory usage percentage
}

// NewWorkerSpawnerWithMemoryLimit creates a new QueuedWorkPoolWithMemoryLimit with the specified memory limit.
func NewWorkerSpawnerWithMemoryLimit(maxMemoryPct uint64) *WorkerSpawnerWithMemoryLimit {
	if maxMemoryPct == 0 {
		// default memory limit is 75%
		maxMemoryPct = 75
	}

	LogInfo("Initializing worker spawner with memory limit...")
	LogInfo(fmt.Sprintf("Memory limit: %d%%", maxMemoryPct))

	return &WorkerSpawnerWithMemoryLimit{
		taskCh:       make(chan Task),
		stopCh:       make(chan struct{}),
		maxMemoryPct: maxMemoryPct,
	}
}

// Start initializes and starts the worker pool.
func (wp *WorkerSpawnerWithMemoryLimit) Start() {
	for {
		select {
		case task, ok := <-wp.taskCh:
			if !ok {
				return // taskCh is closed
			}
			wp.wg.Add(1)
			go func(t Task) {
				defer wp.wg.Done()
				defer func() {
					if r := recover(); r != nil {
						LogWarning(fmt.Sprintf("Recovered from panic in worker goroutine: %v", r))
					}
				}()
				t.Execute()
			}(task)
		case <-wp.stopCh:
			return // stop the goroutine when signaled
		}
	}
}

// Run submits a task to the worker pool.
func (wp *WorkerSpawnerWithMemoryLimit) Run(task Task) (bool, error) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Calculate memory usage percentage
	usedMemoryPct := (memStats.Alloc * 100) / memStats.TotalAlloc

	// If memory usage exceeds the limit, reject the task
	if usedMemoryPct > wp.maxMemoryPct {
		LogWarning(fmt.Sprintf("Memory usage exceeds %d%%. Task rejected.", wp.maxMemoryPct))
		return false, fmt.Errorf("task rejected: memory limit exceeded")
	}

	select {
	case wp.taskCh <- task:
		return true, nil
	default:
		return false, fmt.Errorf("task rejected: queue full")
	}
}

// Wait waits for all goroutines to finish.
func (wp *WorkerSpawnerWithMemoryLimit) Close() {
	close(wp.taskCh)
	wp.wg.Wait()
}
