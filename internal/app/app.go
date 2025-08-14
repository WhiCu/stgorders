package app

import (
	"context"
	"log"
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
	// Create logger
	log := getLogger(&cfg.Logger)
	log.Info("logger created", slog.String("level", cfg.Logger.Level), slog.String("path", cfg.Logger.Path), slog.Int("size", cfg.Logger.Size))

	// Create handler
	h := kc.NewKafkaConsumer(log.With(slog.String("handler", "kafka-consumer")), cfg.Kafka)
	log.Info("handler created", slog.String("brokers", strings.Join(cfg.Kafka.Brokers, ", ")), slog.String("group_id", cfg.Kafka.GroupID), slog.String("topic", cfg.Kafka.Topic))

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
