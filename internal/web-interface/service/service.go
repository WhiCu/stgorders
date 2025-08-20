package service

import (
	"context"
	"log/slog"

	"github.com/WhiCu/stgorders/db/model"
)

type Event func(ctx context.Context) error

type Storage interface {
	GetJsonOrderByUID(ctx context.Context, orderUID string) (jsonOrder *model.JsonOrder, err error)
}

type Service struct {
	storage Storage
	log     *slog.Logger
}

func NewService(storage Storage, log *slog.Logger) *Service {
	return &Service{
		storage: storage,
		log:     log,
	}
}

func (s *Service) GetJsonOrderByUID(ctx context.Context, orderUID string) (jsonOrder *model.JsonOrder, err error) {
	jsonOrder, err = s.storage.GetJsonOrderByUID(ctx, orderUID)
	if err != nil {
		s.log.Error("could not get order", slog.String("ERR", err.Error()))
		return nil, err
	}
	return jsonOrder, nil
}

// func (s *Service) Serve(ctx context.Context, data []byte) (err error) {
// 	db, rollback, commit, err := s.storage.WithTx(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	defer rollback(ctx)

// 	var jo jsonOrder
// 	if err = json.Unmarshal(data, &jo); err != nil {
// 		return err
// 	}
// 	log := s.log.With(slog.String("order_uid", jo.OrderUID))
// 	log.Info("order received")

// 	_, err = db.CreateOrder(ctx, jo.Order)
// 	if err != nil {
// 		log.Error("could not create order", slog.String("ERR", err.Error()))
// 		return err
// 	}

// 	_, err = db.CreateDelivery(ctx, jo.OrderUID, jo.Delivery)
// 	if err != nil {
// 		log.Error("could not create delivery", slog.String("ERR", err.Error()))
// 		return err
// 	}

// 	_, err = db.CreatePayment(ctx, jo.Payment)
// 	if err != nil {
// 		log.Error("could not create payment", slog.String("ERR", err.Error()))
// 		return err
// 	}

// 	for _, it := range jo.Items {
// 		_, err = db.CreateItem(ctx, it)
// 		if err != nil {
// 			s.log.Error("could not create item", slog.String("ERR", err.Error()))
// 			return err
// 		}
// 	}

// 	if err = commit(ctx); err != nil {
// 		return err
// 	}

// 	s.log.Info("order processed", slog.String("order_uid", jo.OrderUID))
// 	return nil
// }
