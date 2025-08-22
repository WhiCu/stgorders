package worker

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	slogmock "github.com/samber/slog-mock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestWorkerPool_Panic(t *testing.T) {
	Convey("Worker recovers from panic and continues processing", t, func() {
		once := sync.Once{}
		var testInt atomic.Int32
		wp := NewWorkerPool(1, func(task int) error {
			once.Do(func() { panic("test panic") })
			testInt.Add(1)
			return nil
		}, 1)

		So(wp.Serve(0), ShouldBeTrue)
		So(testInt.Load(), ShouldEqual, 0)
		So(wp.Serve(0), ShouldBeTrue)
		wp.StopAndWait()
		So(testInt.Load(), ShouldEqual, 1)
	})
}

func TestWorkerPool_Logger(t *testing.T) {
	Convey("Logger is called during processing", t, func() {
		work := false
		wp := NewWorkerPool(1, func(task int) error { return nil }, 1)
		wp.log = slog.New(slogmock.Option{
			Handle: func(ctx context.Context, r slog.Record) error {
				work = true
				return nil
			},
		}.NewMockHandler())

		So(wp.Serve(0), ShouldBeTrue)
		wp.StopAndWait()

		So(work, ShouldBeTrue)
	})
}

func TestWorkerPool_BasicExecution(t *testing.T) {
	Convey("Processes all queued tasks", t, func() {
		var testInt atomic.Int32

		wp := NewWorkerPool(10, func(task int) error {
			testInt.Add(1)
			return nil
		}, 10)

		for i := 0; i < 1000; i++ {
			So(wp.Serve(i), ShouldBeTrue)
		}

		wp.StopAndWait()
		So(testInt.Load(), ShouldEqual, 1000)
	})
}

func TestWorkerPool_ServeAfterStop(t *testing.T) {
	Convey("Serve returns false after StopAndWait", t, func() {
		wp := NewWorkerPool(1, func(task int) error { return nil }, 1)
		wp.StopAndWait()
		So(wp.Serve(0), ShouldBeFalse)
	})
}

func TestWorkerPool_StopAndWaitContext_Timeout(t *testing.T) {
	Convey("StopAndWaitContext returns DeadlineExceeded on timeout", t, func() {
		wp := NewWorkerPool(1, func(task int) error {
			time.Sleep(200 * time.Millisecond)
			return nil
		}, 1)

		So(wp.Serve(0), ShouldBeTrue)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := wp.StopAndWaitContext(ctx)
		So(errors.Is(err, context.DeadlineExceeded), ShouldBeTrue)
	})
}

func TestWorkerPool_StopAndWaitContext_NoTimeout(t *testing.T) {
	Convey("StopAndWaitContext waits successfully without timeout", t, func() {
		wp := NewWorkerPool(1, func(task int) error {
			time.Sleep(10 * time.Millisecond)
			return nil
		}, 1)

		So(wp.Serve(0), ShouldBeTrue)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		err := wp.StopAndWaitContext(ctx)
		So(err, ShouldBeNil)
	})
}
