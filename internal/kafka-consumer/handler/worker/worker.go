package worker

import (
	"context"
	"log/slog"
	"sync"

	"github.com/WhiCu/stgorders/pkg/logger"
)

type WorkerPool[T any] struct {
	wg       sync.WaitGroup
	tasks    chan T
	workerFn func(T) error
	lock     sync.Mutex
	closed   bool

	log *slog.Logger
}

func NewWorkerPool[T any](workers int, workerFn func(T) error, buf int, log ...*slog.Logger) *WorkerPool[T] {

	wp := &WorkerPool[T]{
		tasks:    make(chan T, buf),
		workerFn: workerFn,
	}

	if len(log) == 0 {
		wp.log = slog.New(logger.NewNOPHandler())
	} else {
		wp.log = log[0]
	}

	wp.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			var err error
			log := wp.log.With(slog.Int("workerID", i))
			log.Debug("worker started")
			defer wp.wg.Done()
			for task := range wp.tasks {
				log.Debug("task received")
				err = func() error { // отдельный блок для recover
					defer func() {
						if r := recover(); r != nil {
							log.WithGroup("panic").Error("recovered panic in worker", slog.Any("ERR", r))
						}
					}()

					return workerFn(task)
				}()
				// err = workerFn(task)
				if err != nil {
					wp.log.Error("task failed", slog.Int("workerID", i), slog.String("ERR", err.Error()))
				}
				wp.log.Debug("task processed", slog.Int("workerID", i))
			}
			log.Debug("worker stopped")
		}()

	}
	wp.log.Debug("worcker pool created", slog.Int("workers", workers), slog.Int("buf", buf))
	return wp
}

func (wp *WorkerPool[T]) AddLogger(log *slog.Logger) {
	wp.log = log
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
