package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/WhiCu/stgorders/db/cache"
	"github.com/WhiCu/stgorders/db/model"
	"github.com/WhiCu/stgorders/db/storage"
	"github.com/WhiCu/stgorders/internal/config"
	kc "github.com/WhiCu/stgorders/internal/kafka-consumer"
	wi "github.com/WhiCu/stgorders/internal/web-interface"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	shutdownTimeout = 10 * time.Second
)

type consumer interface {
	ListenAndServe(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type App struct {
	csm  consumer
	srv  *http.Server
	done chan error

	cfg *config.Config
	log *slog.Logger
}

func (a *App) gracefulShutdown(ctx context.Context, cancelParent context.CancelFunc) {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	cancelParent()
	a.log.Info("shutting down gracefully, press Ctrl+C again to force")
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := a.Shutdown(ctx); err != nil {
		a.log.Error("Server forced to shutdown with error", slog.String("ERR", err.Error()))
		a.done <- err
		return
	}

	a.log.Info("Server successfully shutdown")

	a.done <- nil
}

func (a *App) Shutdown(ctx context.Context) (err error) {
	err = a.csm.Shutdown(ctx)
	if err != nil {
		a.log.Debug("could not shutdown consumer", slog.String("ERR", err.Error()))
		return err
	}

	err = a.srv.Shutdown(ctx)
	if err != nil {
		a.log.Debug("could not shutdown server", slog.String("ERR", err.Error()))
		return err
	}

	return nil
}

func NewApp(cfg *config.Config) *App {
	// Create logger
	log := getLogger(&cfg.Logger)
	log.Info("logger created",
		slog.String("level", cfg.Logger.Level),
		slog.String("path", cfg.Logger.Path),
		slog.Int("size", cfg.Logger.Size),
		slog.Bool("compress", cfg.Logger.Compress),
	)

	// Create storage
	p, err := pgxpool.New(context.Background(), cfg.Storage.DSN())
	if err != nil {
		panic(err)
	}
	stg := storage.NewStorage(p, log.WithGroup("storage"))

	// Create cache
	cache := cache.NewLRUCache[string, model.JsonOrder](cfg.Cache.Size, log.WithGroup("cache"))
	err = stg.InitCache(context.Background(), cache)
	if err != nil {
		panic(err)
	}

	// Create kafka consumer
	csm := kc.NewKafkaConsumer(log.WithGroup("kafka-consumer"), cfg.Kafka, stg, cache)
	log.Info("handler created",
		slog.String("brokers", strings.Join(cfg.Kafka.Brokers, ", ")),
		slog.String("group_id", cfg.Kafka.GroupID),
		slog.String("topic", cfg.Kafka.Topic),
	)

	// Create server
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	wi.RegisterRoutes(router.Group("order"), log.WithGroup("web-interface"), stg, cache)
	srv := http.Server{
		Addr:    cfg.Server.ServerAddr(),
		Handler: router,
	}

	return &App{
		csm:  csm,
		srv:  &srv,
		done: make(chan error),
		cfg:  cfg,
		log:  log,
	}
}

func (a *App) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// graceful shutdown слушает сигналы и инициирует остановку
	go a.gracefulShutdown(ctx, cancel)

	// запускаем HTTP сервер
	go func() {
		a.log.Info("starting http server", slog.String("addr", a.srv.Addr))
		if err := a.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Error("http server failed", slog.String("ERR", err.Error()))
			a.done <- err
		}
	}()

	// запускаем kafka consumer
	go func() {
		a.log.Info("starting kafka consumer")
		if err := a.csm.ListenAndServe(ctx); err != nil {
			a.log.Error("kafka consumer failed", slog.String("ERR", err.Error()))
			a.done <- err
		}
	}()

	// ждём либо ошибку, либо завершение gracefulShutdown
	return <-a.done
}

// func (a *App) Run(ctx context.Context) error {
// 	ctx, cancel := context.WithCancel(ctx)
// 	defer cancel()

// 	go a.gracefulShutdown(ctx, cancel)

// 	a.log.Info("starting server")
// 	if err := a.csm.ListenAndServe(ctx); err != nil {
// 		a.log.Error("could not listen", slog.String("ERR", err.Error()))
// 		return err
// 	}

// 	return <-a.done
// }

// func (a *App) RunWithRecover(ctx context.Context) (err error) {
// 	for {
// 		err = func() error {
// 			defer func() {
// 				if r := recover(); r != nil {
// 					log := a.log.WithGroup("panic")
// 					log.Error("recovered from panic", slog.Any("ERR", r))
// 				}
// 			}()
// 			return a.Run(ctx)
// 		}()
// 		if err != nil {
// 			a.log.Error("could not run", slog.String("ERR", err.Error()))
// 			return err
// 		}
// 	}
// }
