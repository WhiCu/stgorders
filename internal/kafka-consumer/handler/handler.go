package handler

import (
	"context"
	"log/slog"

	"github.com/WhiCu/stgorders/internal/kafka-consumer/handler/worker"
	"github.com/segmentio/kafka-go"
)

type Handler struct {
	// main consumer
	consumer *kafka.Reader

	// worker
	workerPool *worker.WorkerPool[*kafka.Message]

	// logger
	log *slog.Logger
}

type ConsumerConfig struct {
	Brokers []string
	GroupID string
	Topic   string
}

func NewHandler(log *slog.Logger, cfg ConsumerConfig) *Handler {
	c := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers: cfg.Brokers,
			GroupID: cfg.GroupID,
			Topic:   cfg.Topic,
		},
	)

	wp := worker.NewWorkerPool(
		10,
		func(m *kafka.Message) error {
			log.Info("message", "topic", m.Topic, "partition", m.Partition, "offset", m.Offset, "key", string(m.Key), "value", string(m.Value))
			return c.CommitMessages(context.Background(), *m)
		},
		10,
	)

	return &Handler{
		consumer:   c,
		workerPool: wp,
		log:        log,
	}
}

func (h *Handler) ListenAndServe(ctx context.Context) error {

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			h.log.Info("fetching message")
			m, err := h.consumer.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return nil
				}
				return err
			}
			h.log.Info("message fetched")
			if !h.workerPool.Serve(&m) {
				return nil
			}
			h.log.Info("message processed")
		}
	}
}

func (h *Handler) Shutdown(ctx context.Context) (err error) {
	h.log.Info("shutting down")
	if err = h.consumer.Close(); err != nil {
		return err
	}
	h.log.Info("shutting down workers")
	if err = h.workerPool.StopAndWaitContext(ctx); err != nil {
		return err
	}
	h.log.Info("shutting down done")
	return nil
}
func (h *Handler) Close() {
	h.consumer.Close()
}
