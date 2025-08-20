package client

import (
	"context"
	"log/slog"

	"github.com/WhiCu/stgorders/db/cache"
	"github.com/WhiCu/stgorders/db/model"
	"github.com/WhiCu/stgorders/db/storage"
)

type StorageAdapter struct {
	*storage.Storage
	cache *cache.LRUCache[string, model.JsonOrder]
	log   *slog.Logger
}

func NewStorage(storage *storage.Storage, cache *cache.LRUCache[string, model.JsonOrder], log *slog.Logger) *StorageAdapter {
	log.Debug("creating storage adapter")
	return &StorageAdapter{
		Storage: storage,
		cache:   cache,
		log:     log,
	}
}

func (s *StorageAdapter) GetJsonOrderByUID(ctx context.Context, orderUID string) (jsonOrder *model.JsonOrder, err error) {
	log := s.log.With(slog.String("order_uid", orderUID))

	if jo, err := s.cache.Get(orderUID); err == nil {
		log.Debug("order found in cache")
		return &jo, nil
	}

	order, err := s.GetOrderByUID(ctx, orderUID)
	if err != nil {
		log.Error("could not get order", slog.String("ERR", err.Error()))
		return nil, err
	}
	delivery, err := s.GetDeliveryByOrderUID(ctx, order.OrderUid)
	if err != nil {
		log.Error("could not get delivery", slog.String("ERR", err.Error()))
		return nil, err
	}
	payment, err := s.GetPaymentByTransaction(ctx, order.OrderUid)
	if err != nil {
		log.Error("could not get payment", slog.String("ERR", err.Error()))
		return nil, err
	}
	items, err := s.GetItemsByTrackNumber(ctx, order.TrackNumber)
	if err != nil {
		log.Error("could not get items", slog.String("ERR", err.Error()))
		return nil, err
	}

	jo := model.ModelJSONOrder(order, delivery, payment, items)

	if err = s.cache.Set(orderUID, jo); err != nil {
		log.Error("could not set order in cache", slog.String("ERR", err.Error()))
		return nil, err
	}

	return &jo, nil
}
