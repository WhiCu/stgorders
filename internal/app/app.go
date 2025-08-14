package app

import (
	"context"
	"log/slog"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/WhiCu/stgorders/internal/config"
	kc "github.com/WhiCu/stgorders/internal/kafka-consumer"
)

type consumer interface {
	ListenAndServe(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type App struct {
	server consumer
	done   chan error

	cfg *config.Config
	log *slog.Logger
}

func (a *App) gracefulShutdown(cl context.CancelFunc) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	cl()
	a.log.Info("shutting down gracefully, press Ctrl+C again to force")
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		a.log.Error("Server forced to shutdown with error", slog.String("ERR", err.Error()))
		a.done <- err
		return
	}

	a.log.Info("Server successfully shutdown")

	a.done <- nil
}

func NewApp(cfg *config.Config) *App {
	// Create logger
	log := getLogger(&cfg.Logger)
	log.Info("logger created", slog.String("level", cfg.Logger.Level), slog.String("path", cfg.Logger.Path), slog.Int("size", cfg.Logger.Size))

	// Create handler
	h := kc.NewKafkaConsumer(log.WithGroup("kafka-consumer"), cfg.Kafka)
	log.Info("handler created", slog.String("brokers", strings.Join(cfg.Kafka.Brokers, ", ")), slog.String("group_id", cfg.Kafka.GroupID), slog.String("topic", cfg.Kafka.Topic))

	return &App{
		server: h,
		done:   make(chan error),
		cfg:    cfg,
		log:    log,
	}
}

func (a *App) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go a.gracefulShutdown(cancel)

	if err := a.server.ListenAndServe(ctx); err != nil && err != http.ErrServerClosed {
		a.log.Error("could not listen", slog.String("ERR", err.Error()))
		return err
	}

	return <-a.done
}
