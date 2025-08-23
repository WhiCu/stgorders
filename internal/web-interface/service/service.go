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
	log := s.log.With(slog.String("order_uid", orderUID))
	jsonOrder, err = s.storage.GetJsonOrderByUID(ctx, orderUID)
	if err != nil {
		log.Error("could not get order", slog.String("ERR", err.Error()))
		return nil, err
	}
	log.Info("order found")
	return jsonOrder, nil
}
