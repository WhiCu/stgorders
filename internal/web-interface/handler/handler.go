package handler

import (
	"context"
	"log/slog"

	"github.com/WhiCu/stgorders/db/model"
)

type service interface {
	GetJsonOrderByUID(ctx context.Context, orderUID string) (jsonOrder *model.JsonOrder, err error)
}

type Handler struct {
	service service
	log     *slog.Logger
}

func NewHandler(service service, log *slog.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}
