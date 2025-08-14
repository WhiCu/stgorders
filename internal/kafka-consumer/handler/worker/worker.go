package worker

import (
	"context"
	"sync"
)

type WorkerPool[T any] struct {
	wg       sync.WaitGroup
	tasks    chan T
	workerFn func(T) error
	lock     sync.Mutex
	closed   bool
}

func NewWorkerPool[T any](workers int, workerFn func(T) error, buf int) *WorkerPool[T] {
	wp := &WorkerPool[T]{
		tasks:    make(chan T, buf),
		workerFn: workerFn,
	}

	wp.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wp.wg.Done()
			for task := range wp.tasks {
				_ = workerFn(task)
			}
		}()
	}

	return wp
}

func (wp *WorkerPool[T]) Serve(task T) bool {
	wp.lock.Lock()
	defer wp.lock.Unlock()
	if wp.closed {
		return false
	}
	wp.tasks <- task
	return true
}

func (wp *WorkerPool[T]) StopAndWait() {
	wp.lock.Lock()
	if wp.closed {
		wp.lock.Unlock()
		return
	}
	wp.closed = true
	close(wp.tasks)
	wp.lock.Unlock()

	wp.wg.Wait()
}

func (wp *WorkerPool[T]) StopAndWaitContext(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		wp.StopAndWait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
