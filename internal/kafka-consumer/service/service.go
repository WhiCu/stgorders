package service

import (
	"log/slog"
)

type Storage interface {
}

type Service struct {
	storage Storage

	log *slog.Logger
}

func NewService(storage Storage, log *slog.Logger) *Service {
	return &Service{
		storage: storage,
		log:     log,
	}
}

func (s *Service) Serve(data []byte) error {
	s.log.Info("message processed", slog.String("data", string(data)))
	return nil
}
