package worker

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorkerPool_BasicExecution(t *testing.T) {
	var testInt atomic.Int32

	wp := NewWorkerPool(10, func(task int) error {
		testInt.Add(1)
		return nil
	}, 10)

	for i := 0; i < 1000; i++ {
		if !wp.Serve(i) {
			t.Fatalf("Serve returned false unexpectedly for task %d", i)
		}
	}

	wp.StopAndWait()
	assert.Equal(t, int32(1000), testInt.Load(), "expected 1000, got %d", testInt.Load())
}

func TestWorkerPool_ServeAfterStop(t *testing.T) {
	wp := NewWorkerPool(1, func(task int) error { return nil }, 1)

	wp.StopAndWait()

	if wp.Serve(42) {
		t.Fatalf("Serve should return false after StopAndWait")
	}
}

func TestWorkerPool_StopAndWaitContext_Timeout(t *testing.T) {
	wp := NewWorkerPool(1, func(task int) error {
		time.Sleep(200 * time.Millisecond)
		return nil
	}, 1)

	// Задача, которая будет выполняться долго
	wp.Serve(1)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := wp.StopAndWaitContext(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
}

func TestWorkerPool_StopAndWaitContext_NoTimeout(t *testing.T) {
	wp := NewWorkerPool(1, func(task int) error {
		time.Sleep(10 * time.Millisecond)
		return nil
	}, 1)

	wp.Serve(1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := wp.StopAndWaitContext(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
