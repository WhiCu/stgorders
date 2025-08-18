package kc

import (
	"log/slog"

	"github.com/WhiCu/stgorders/internal/config"
	"github.com/WhiCu/stgorders/internal/kafka-consumer/client"
	"github.com/WhiCu/stgorders/internal/kafka-consumer/handler"
	"github.com/WhiCu/stgorders/internal/kafka-consumer/service"
)

func NewKafkaConsumer(log *slog.Logger, cfg config.KafkaConfig, storage *client.Storage) *handler.Handler {
	return handler.NewHandler(
		log.WithGroup("handler"),
		handler.ConsumerConfig{
			Brokers:        cfg.Brokers,
			GroupID:        cfg.GroupID,
			Topic:          cfg.Topic,
			WorkerPoolSize: cfg.WorkerPool.Size,
			WorkerPoolBuf:  cfg.WorkerPool.Buf,
		},
		initService(log, storage),
	)
}

func initService(log *slog.Logger, storage *client.Storage) *service.Service {
	return service.NewService(storage, log.WithGroup("service"))
}
