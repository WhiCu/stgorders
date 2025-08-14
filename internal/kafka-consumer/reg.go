package kc

import (
	"log/slog"

	"github.com/WhiCu/stgorders/internal/config"
	"github.com/WhiCu/stgorders/internal/kafka-consumer/handler"
	"github.com/WhiCu/stgorders/internal/kafka-consumer/service"
)

func NewKafkaConsumer(log *slog.Logger, cfg config.KafkaConfig) *handler.Handler {
	return handler.NewHandler(log.WithGroup("handler"), handler.ConsumerConfig{
		Brokers:        cfg.Brokers,
		GroupID:        cfg.GroupID,
		Topic:          cfg.Topic,
		WorkerPoolSize: cfg.WorkerPool.Size,
		WorkerPoolBuf:  cfg.WorkerPool.Buf,
	},
		initService(log),
	)
}

func initService(log *slog.Logger) *service.Service {
	return service.NewService(nil, log.WithGroup("service"))
}
