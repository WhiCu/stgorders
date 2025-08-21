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
	"golang.org/x/sync/errgroup"
)

var (
	shutdownTimeout = 10 * time.Second
)

type consumer interface {
	ListenAndServe(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type App struct {
	csm consumer
	srv *http.Server

	cfg *config.Config
	log *slog.Logger
}

func (a *App) gracefulShutdown(ctx context.Context, cancelParent context.CancelFunc) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	a.log.Info("shutting down gracefully, press Ctrl+C again to force")
	stop()
	cancelParent()

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := a.Shutdown(ctx); err != nil {
		a.log.Error("Server forced to shutdown with error", slog.String("ERR", err.Error()))
		return err
	}

	a.log.Info("Server successfully shutdown")

	return nil
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
	if err = p.Ping(context.Background()); err != nil {
		panic(err)
	}

	stg := storage.NewStorage(p, log.WithGroup("storage"))

	// Create cache
	log.Info("cache created", slog.Int("size", cfg.Cache.Size))
	var c cache.Cache[string, model.JsonOrder]
	c = cache.NewNOPCache[string](model.JsonOrder{})
	if cfg.Cache.Size != 0 {
		log.Info("LRU cache created", slog.Int("size", cfg.Cache.Size))
		c = cache.NewLRUCache[string, model.JsonOrder](cfg.Cache.Size, log.WithGroup("cache"))
	}

	err = stg.InitCache(context.Background(), c)
	if err != nil {
		panic(err)
	}

	// Create kafka consumer
	csm := kc.NewKafkaConsumer(log.WithGroup("kafka-consumer"), cfg.Kafka, stg, c)
	log.Info("handler created",
		slog.String("brokers", strings.Join(cfg.Kafka.Brokers, ", ")),
		slog.String("group_id", cfg.Kafka.GroupID),
		slog.String("topic", cfg.Kafka.Topic),
	)

	// Create server
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	wi.RegisterRoutes(router.Group("order"), log.WithGroup("web-interface"), stg, c)
	srv := http.Server{
		Addr:    cfg.Server.ServerAddr(),
		Handler: router,
	}

	return &App{
		csm: csm,
		srv: &srv,
		cfg: cfg,
		log: log,
	}
}

func (a *App) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		a.log.Info("starting http server", slog.String("addr", a.srv.Addr))
		if err := a.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Error("http server failed", slog.String("ERR", err.Error()))
			return err
		}
		return nil
	})

	eg.Go(func() error {
		a.log.Info("starting kafka consumer", slog.String("brokers", strings.Join(a.cfg.Kafka.Brokers, ", ")))
		if err := a.csm.ListenAndServe(ctx); err != nil && !errors.Is(err, context.Canceled) {
			a.log.Error("kafka consumer failed", slog.String("ERR", err.Error()))
			return err
		}
		return nil
	})

	eg.Go(func() error {
		return a.gracefulShutdown(ctx, cancel)
	})

	return eg.Wait()
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
