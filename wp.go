package workerpool

import "sync/atomic"

type WP struct {
	workers        int64
	taskQueue      chan func()
	stopController chan struct{}
}

func New() *WP {
	wp := &WP{
		taskQueue:      make(chan func()),
		stopController: make(chan struct{}),
	}
	return wp
}

func (wp *WP) Submit(task func()) {
	if task != nil {
		wp.taskQueue <- task
	}
}

func (wp *WP) Stop() {

}

func (wp *WP) Dispatch(workers int) {
	for i := 1; i < workers; i++ {
		atomic.AddInt64(&wp.workers, 1)
		go func() {
		LOOP:
			for {
				select {
				case <-wp.stopController:
					break LOOP
				default:
					select {
					case task, ok := <-wp.taskQueue:
						if !ok {
							break LOOP
						}
						if task != nil {
							task()
						}
					default:
					}
				}
			}
		}()
	}
}
