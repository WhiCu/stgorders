package handler

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/WhiCu/stgorders/internal/kafka-consumer/handler/worker"
	"github.com/WhiCu/stgorders/pkg/logger"
	"github.com/segmentio/kafka-go"
	. "github.com/smartystreets/goconvey/convey"
)

type mockConsumer struct {
	errCommit error
	errClose  error
	errFetch  error
	messages  []kafka.Message
	msgIndex  int
}

func (m *mockConsumer) Close() error { return m.errClose }

func (m *mockConsumer) CommitMessages(ctx context.Context, messages ...kafka.Message) error {
	return m.errCommit
}

var errNoMoreMessages = errors.New("no more messages")

func (m *mockConsumer) FetchMessage(ctx context.Context) (kafka.Message, error) {
	if m.errFetch != nil {
		return kafka.Message{}, m.errFetch
	}
	if m.msgIndex < len(m.messages) {
		msg := m.messages[m.msgIndex]
		m.msgIndex++
		return msg, nil
	}
	return kafka.Message{}, errNoMoreMessages
}

type mockService struct {
	err error
}

func (s *mockService) Serve(_ context.Context, _ []byte) error { return s.err }
func (s *mockService) Close() error                            { return nil }

type mockWorkerPool struct {
	err        error
	serveCalls int
	maxCalls   int
	countCalls int
}

func (m *mockWorkerPool) Serve(_ *kafka.Message) bool {
	m.serveCalls++
	if m.maxCalls != 0 && m.serveCalls > m.maxCalls {
		return false
	}
	m.countCalls++
	return true
}

func (m *mockWorkerPool) StopAndWaitContext(_ context.Context) error { return m.err }

func GetWorkerPool() *worker.WorkerPool[*kafka.Message] {
	wp := worker.NewWorkerPool(
		1,
		func(m *kafka.Message) (err error) {
			return nil
		},
		1,
		logger.NewNOPSlog(),
	)
	return wp
}

func TestHandler_ListenAndServe(t *testing.T) {
	Convey("ListenAndServe stops on context done", t, func() {
		mc := &mockConsumer{}
		sv := &mockService{}
		wp := &mockWorkerPool{}
		h := newHandler(logger.NewNOPSlog(), mc, sv, wp)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := h.ListenAndServe(ctx)
		So(err, ShouldBeNil)
	})

	Convey("ListenAndServe returns fetch error", t, func() {
		mc := &mockConsumer{errFetch: errors.New("fetch failed")}
		sv := &mockService{}
		wp := &mockWorkerPool{}
		h := newHandler(logger.NewNOPSlog(), mc, sv, wp)
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := h.ListenAndServe(ctx)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "fetch failed")
	})

	Convey("ListenAndServe continues on EOF error", t, func() {
		mc := &mockConsumer{errFetch: io.EOF}
		sv := &mockService{}
		wp := &mockWorkerPool{}
		h := newHandler(logger.NewNOPSlog(), mc, sv, wp)
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := h.ListenAndServe(ctx)
		So(err, ShouldBeNil)
	})

	Convey("ListenAndServe stops when worker pool is full", t, func() {
		mc := &mockConsumer{
			messages: []kafka.Message{
				{Value: []byte("message1")},
				{Value: []byte("message2")},
			},
		}
		sv := &mockService{}
		wp := &mockWorkerPool{maxCalls: 1}
		h := newHandler(logger.NewNOPSlog(), mc, sv, wp)
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := h.ListenAndServe(ctx)
		So(err, ShouldBeNil)
		So(wp.countCalls, ShouldEqual, 1)
	})

	Convey("ListenAndServe returns nil if context is cancelled during fetch", t, func() {
		mc := &mockConsumer{errFetch: context.Canceled}
		sv := &mockService{}
		wp := &mockWorkerPool{}
		h := newHandler(logger.NewNOPSlog(), mc, sv, wp)
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := h.ListenAndServe(ctx)
		So(err, ShouldBeNil)
	})

	Convey("ListenAndServe continues on service error and stops on context timeout", t, func() {
		msg := kafka.Message{Value: []byte("bad message")}
		mc := &mockConsumer{messages: []kafka.Message{msg}}
		sv := &mockService{err: errors.New("service error")}
		wp := &mockWorkerPool{}
		h := newHandler(logger.NewNOPSlog(), mc, sv, wp)
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := h.ListenAndServe(ctx)
		So(err, ShouldEqual, errNoMoreMessages)
	})
}

func TestHandler_Shutdown(t *testing.T) {
	Convey("Shutdown succeeds with normal worker pool and consumer", t, func() {
		mc := &mockConsumer{}
		sv := &mockService{}
		wp := &mockWorkerPool{}
		h := newHandler(logger.NewNOPSlog(), mc, sv, wp)
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := h.Shutdown(ctx)
		So(err, ShouldBeNil)
	})

	Convey("Shutdown returns error when worker pool fails", t, func() {
		mc := &mockConsumer{}
		sv := &mockService{}
		wp := &mockWorkerPool{err: errors.New("worker pool stop failed")}
		h := newHandler(logger.NewNOPSlog(), mc, sv, wp)
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := h.Shutdown(ctx)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "worker pool stop failed")
	})

	Convey("Shutdown returns error when consumer fails to close", t, func() {
		mc := &mockConsumer{errClose: errors.New("consumer close failed")}
		sv := &mockService{}
		wp := &mockWorkerPool{}
		h := newHandler(logger.NewNOPSlog(), mc, sv, wp)
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := h.Shutdown(ctx)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "consumer close failed")
	})
}

func TestHandler_Close(t *testing.T) {
	Convey("Close succeeds when consumer closes without error", t, func() {
		mc := &mockConsumer{}
		sv := &mockService{}
		wp := &mockWorkerPool{}
		h := newHandler(logger.NewNOPSlog(), mc, sv, wp)

		err := h.Close()
		So(err, ShouldBeNil)
	})

	Convey("Close returns error when consumer fails to close", t, func() {
		mc := &mockConsumer{errClose: errors.New("close failed")}
		sv := &mockService{}
		wp := &mockWorkerPool{}
		h := newHandler(logger.NewNOPSlog(), mc, sv, wp)

		err := h.Close()
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "close failed")
	})
}

func TestNewHandler(t *testing.T) {
	Convey("NewHandler returns valid handler", t, func() {
		cfg := ConsumerConfig{
			Brokers:        []string{"localhost:9092"},
			GroupID:        "test-group",
			Topic:          "test-topic",
			WorkerPoolSize: 2,
			WorkerPoolBuf:  10,
		}
		sv := &mockService{}

		h := NewHandler(logger.NewNOPSlog(), cfg, sv)

		So(h, ShouldNotBeNil)
		So(h.consumer, ShouldNotBeNil)
		So(h.workerPool, ShouldNotBeNil)
		So(h.service, ShouldEqual, sv)
		So(h.log, ShouldNotBeNil)
	})
}
