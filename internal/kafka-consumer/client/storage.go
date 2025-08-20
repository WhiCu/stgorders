package client

import (
	"context"
	"log/slog"

	"github.com/WhiCu/stgorders/db/cache"
	"github.com/WhiCu/stgorders/db/model"
	"github.com/WhiCu/stgorders/db/pg"
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

// func (s *StorageAdapter) WithTx(ctx context.Context) (storage service.Storage, rollback service.Event, commit service.Event, err error) {
// 	stg, rb, cmt, err := s.Storage.WithTx(ctx)
// 	return &StorageAdapter{stg, s.log}, service.Event(rb), service.Event(cmt), err
// }

func (s *StorageAdapter) Save(ctx context.Context, order model.JsonOrder) (err error) {

	storage, rollback, commit, err := s.WithTx(ctx)
	if err != nil {
		s.log.Error("could not create transaction", slog.String("ERR", err.Error()))
		return err
	}
	defer rollback(ctx)

	_, err = CreateOrder(ctx, storage, order.Order)
	if err != nil {
		s.log.Error("could not create order", slog.String("ERR", err.Error()))
		return err
	}

	_, err = CreateDelivery(ctx, storage, order.OrderUID, order.Delivery)
	if err != nil {
		s.log.Error("could not create delivery", slog.String("ERR", err.Error()))
		return err
	}

	_, err = CreatePayment(ctx, storage, order.Payment)
	if err != nil {
		s.log.Error("could not create payment", slog.String("ERR", err.Error()))
		return err
	}

	for _, item := range order.Items {
		_, err = CreateItem(ctx, storage, item)
		if err != nil {
			s.log.Error("could not create item", slog.String("ERR", err.Error()))
			return err
		}
	}

	if err = commit(ctx); err != nil {
		s.log.Error("could not commit transaction", slog.String("ERR", err.Error()))
		return err
	}

	if err = s.cache.Set(order.OrderUID, order); err != nil {
		s.log.Error("could not set cache", slog.String("ERR", err.Error()))
	}
	return nil
}

func (s *StorageAdapter) Close() error {
	s.log.Debug("closing storage adapter")
	s.Storage.Close()
	return nil
}

func CreateOrder(ctx context.Context, storage *storage.Storage, params model.Order) (int64, error) {
	return storage.CreateOrder(ctx, pg.CreateOrderParams{
		OrderUid:          params.OrderUID,
		TrackNumber:       params.TrackNumber,
		Entry:             params.Entry,
		Locale:            params.Locale,
		InternalSignature: params.InternalSignature,
		CustomerID:        params.CustomerID,
		DeliveryService:   params.DeliveryService,
		Shardkey:          params.Shardkey,
		SmID:              params.SmID,
		DateCreated:       params.DateCreated,
		OofShard:          params.OofShard,
	})
}

func CreateDelivery(ctx context.Context, storage *storage.Storage, orderID string, params model.Delivery) (int64, error) {
	return storage.CreateDelivery(ctx, pg.CreateDeliveryParams{
		OrderID: orderID,
		Name:    params.Name,
		Phone:   params.Phone,
		Zip:     params.Zip,
		City:    params.City,
		Address: params.Address,
		Region:  params.Region,
		Email:   params.Email,
	})
}

func CreatePayment(ctx context.Context, storage *storage.Storage, params model.Payment) (int64, error) {
	return storage.CreatePayment(ctx, pg.CreatePaymentParams{
		Transaction:  params.Transaction,
		RequestID:    params.RequestID,
		Currency:     params.Currency,
		Provider:     params.Provider,
		Amount:       params.Amount,
		PaymentDt:    params.PaymentDt,
		Bank:         params.Bank,
		DeliveryCost: params.DeliveryCost,
		GoodsTotal:   params.GoodsTotal,
		CustomFee:    params.CustomFee,
	})
}

func CreateItem(ctx context.Context, storage *storage.Storage, params model.Item) (int64, error) {
	return storage.CreateItem(ctx, pg.CreateItemParams{
		ChrtID:      params.ChrtID,
		TrackNumber: params.TrackNumber,
		Price:       params.Price,
		Rid:         params.Rid,
		Name:        params.Name,
		Sale:        params.Sale,
		Size:        params.Size,
		TotalPrice:  params.TotalPrice,
		NmID:        params.NmID,
		Brand:       params.Brand,
		Status:      params.Status,
	})
}
