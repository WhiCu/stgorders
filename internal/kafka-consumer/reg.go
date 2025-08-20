package kc

import (
	"log/slog"

	"github.com/WhiCu/stgorders/db/cache"
	"github.com/WhiCu/stgorders/db/model"
	"github.com/WhiCu/stgorders/db/storage"
	"github.com/WhiCu/stgorders/internal/config"
	"github.com/WhiCu/stgorders/internal/kafka-consumer/client"
	"github.com/WhiCu/stgorders/internal/kafka-consumer/handler"
	"github.com/WhiCu/stgorders/internal/kafka-consumer/service"
)

func NewKafkaConsumer(log *slog.Logger, cfg config.KafkaConfig, storage *storage.Storage, cache *cache.LRUCache[string, model.JsonOrder]) *handler.Handler {
	stg := client.NewStorage(storage, cache, log.WithGroup("storageAdapter"))
	srv := service.NewService(stg, log.WithGroup("service"))
	return handler.NewHandler(
		log.WithGroup("handler"),
		handler.ConsumerConfig{
			Brokers:        cfg.Brokers,
			GroupID:        cfg.GroupID,
			Topic:          cfg.Topic,
			WorkerPoolSize: cfg.WorkerPool.Size,
			WorkerPoolBuf:  cfg.WorkerPool.Buf,
		},
		srv,
	)
}
