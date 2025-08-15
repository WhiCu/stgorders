package handler

import (
	"context"
	"log/slog"
	"strings"

	"github.com/WhiCu/stgorders/internal/kafka-consumer/handler/worker"
	"github.com/segmentio/kafka-go"
)

type service interface {
	Serve(ctx context.Context, data []byte) error
}

type Handler struct {
	// main consumer
	consumer *kafka.Reader

	// worker
	workerPool *worker.WorkerPool[*kafka.Message]

	// service
	service service

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
		func(m *kafka.Message) (err error) {
			log.Debug("message processed", slog.String("topic", m.Topic), slog.Time("partition", m.Time))
			s.Serve(context.Background(), m.Value)
			if err = c.CommitMessages(context.Background(), *m); err != nil {
				log.Error("could not commit message", slog.String("ERR", err.Error()))
			}
			return err
		},
		cfg.WorkerPoolBuf,
	)

	return &Handler{
		consumer:   c,
		workerPool: wp,
		log:        log,
		service:    s,
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
					h.log.Debug("context canceled", slog.String("ERR", err.Error()))
					return nil
				}
				h.log.Error("could not fetch message", slog.String("ERR", err.Error()))
				return err
			}
			h.log.Debug("message fetched")
			if !h.workerPool.Serve(&m) {
				h.log.Debug("WorkerPool is full")
				return nil
			}
			h.log.Debug("message processed")
		}
	}
}

func (h *Handler) Shutdown(ctx context.Context) (err error) {
	h.log.Debug("shutting down workers")
	if err = h.workerPool.StopAndWaitContext(ctx); err != nil {
		h.log.Error("could not stop workers", slog.String("ERR", err.Error()))
		return err
	}
	h.log.Debug("shutting down")
	if err = h.Close(); err != nil {
		h.log.Error("could not close consumer", slog.String("ERR", err.Error()))
		return err
	}
	h.log.Debug("shutting down done")
	return nil
}
func (h *Handler) Close() error {
	return h.consumer.Close()
}
