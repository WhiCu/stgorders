package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/WhiCu/stgorders/internal/kafka-consumer/handler"
)

func main() {
	cfg := handler.ConsumerConfig{
		Brokers: []string{"localhost:9092"},
		GroupID: "my-group",
		Topic:   "test-topic",
	}

	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	s := http.Server{}
	s.ListenAndServe()
	h := handler.NewHandler(log, cfg)
	if err := h.ListenAndServe(context.Background()); err != nil {
		log.Error("failed to serve and listen", slog.String("ERR", err.Error()))
	}
}
