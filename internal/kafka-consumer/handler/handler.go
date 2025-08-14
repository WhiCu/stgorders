package handler

import (
	"context"
	"log/slog"
	"strings"

	"github.com/WhiCu/stgorders/internal/kafka-consumer/handler/worker"
	"github.com/segmentio/kafka-go"
)

type service interface {
	Serve(data []byte) error
}

type Handler struct {
	// main consumer
	consumer *kafka.Reader

	// worker
	workerPool *worker.WorkerPool[*kafka.Message]

	// service
	service *service

	// logger
	log *slog.Logger
}

type ConsumerConfig struct {
	// Reader Config
	Brokers []string
	GroupID string
	Topic   string

	// Worker Pool Config
	WorkerPoolSize int
	WorkerPoolBuf  int
}

func NewHandler(log *slog.Logger, cfg ConsumerConfig, s service) *Handler {
	kafkacfg := kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		GroupID: cfg.GroupID,
		Topic:   cfg.Topic,
	}
	log.Debug("kafka reader config", slog.String("brokers", strings.Join(kafkacfg.Brokers, ", ")), slog.String("group_id", kafkacfg.GroupID), slog.String("topic", kafkacfg.Topic))
	c := kafka.NewReader(
		kafkacfg,
	)

	wp := worker.NewWorkerPool(
		cfg.WorkerPoolSize,
		func(m *kafka.Message) error {
			log.Debug("message processed", slog.String("topic", m.Topic), slog.Time("partition", m.Time))
			s.Serve(m.Value)
			return c.CommitMessages(context.Background(), *m)
		},
		cfg.WorkerPoolBuf,
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
			h.log.Debug("fetching message")
			m, err := h.consumer.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return nil
				}
				return err
			}
			h.log.Debug("message fetched")
			if !h.workerPool.Serve(&m) {
				return nil
			}
			h.log.Debug("message processed")
		}
	}
}

func (h *Handler) Shutdown(ctx context.Context) (err error) {
	h.log.Debug("shutting down workers")
	if err = h.workerPool.StopAndWaitContext(ctx); err != nil {
		return err
	}
	h.log.Debug("shutting down")
	if err = h.Close(); err != nil {
		return err
	}
	h.log.Debug("shutting down done")
	return nil
}
func (h *Handler) Close() error {
	return h.consumer.Close()
}
