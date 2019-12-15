package workerpool

import "sync/atomic"

// WP worker pool structure
type WP struct {
	workers   int64
	taskQueue chan func()
}

// New creates new worker pool
func New() *WP {
	return &WP{taskQueue: make(chan func())}
}

// Submit put task to process queue
func (wp *WP) Submit(task func()) {
	if task != nil {
		wp.taskQueue <- task
	}
}

// Stop stops all workers
func (wp *WP) Stop() {
	close(wp.taskQueue)
}

// Dispatch runs provided number of workers
func (wp *WP) Dispatch(workers int) {
	for i := 1; i < workers; i++ {
		wp.workers++
		go func() {
			for task := range wp.taskQueue {
				task()
			}
			atomic.AddInt64(&wp.workers, -1)
		}()
	}
}

// Workers returns count of runnable workers
func (wp *WP) Workers() int64 {
	return wp.workers
}
