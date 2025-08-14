package app

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WhiCu/stgorders/internal/config"
	"github.com/WhiCu/stgorders/internal/kafka-consumer/handler"
)

type consumer interface {
	ListenAndServe(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type App struct {
	server consumer
	done   chan error

	cfg *config.Config
}

func (a *App) gracefulShutdown(cl context.CancelFunc) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	cl()
	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
		a.done <- err
		return
	}

	log.Println("Server exiting")

	a.done <- nil
}

func NewApp(cfg *config.Config) *App {

	c := handler.ConsumerConfig{
		Brokers: []string{"localhost:9092"},
		GroupID: "my-group",
		Topic:   "test-topic",
	}

	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	h := handler.NewHandler(log, c)

	return &App{
		server: h,
		done:   make(chan error),
		cfg:    cfg,
	}
}

func (a *App) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go a.gracefulShutdown(cancel)

	if err := a.server.ListenAndServe(ctx); err != nil && err != http.ErrServerClosed {
		log.Printf("could not listen: %s\n", err)
		return err
	}

	return <-a.done
}
