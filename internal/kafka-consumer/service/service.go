package service

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/WhiCu/stgorders/db/model"
	"github.com/go-playground/validator/v10"
)

type Event func(ctx context.Context) error

type Storage interface {
	Save(ctx context.Context, order model.JsonOrder) error

	Close() error
}

type Service struct {
	storage Storage
	log     *slog.Logger
	valid   *validator.Validate
}

func NewService(storage Storage, log *slog.Logger) *Service {
	return &Service{
		storage: storage,
		log:     log,
		valid:   validator.New(),
	}
}

func (s *Service) Serve(ctx context.Context, data []byte) (err error) {
	var jo model.JsonOrder
	if err = json.Unmarshal(data, &jo); err != nil {
		s.log.Error("could not unmarshal data", slog.String("ERR", err.Error()))
		return err
	}
	log := s.log.With(slog.String("order_uid", jo.OrderUID))

	if err = s.valid.Struct(jo); err != nil {
		log.Error("could not validate data", slog.String("ERR", err.Error()))
		return err
	}

	err = s.storage.Save(ctx, jo)
	if err != nil {
		log.Error("could not save data", slog.String("ERR", err.Error()))
		return err
	}
	log.Info("order saved")
	return nil
}

func (s *Service) Close() error {
	return s.storage.Close()
}
