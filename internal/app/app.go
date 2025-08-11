package app

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/WhiCu/stgorders/internal/config"
	"github.com/WhiCu/stgorders/internal/test"
	"github.com/gin-gonic/gin"
)

type App struct {
	server *http.Server
	done   chan error

	cfg *config.Config
}

func (a *App) gracefulShutdown() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

	router := gin.New()

	RegisterMiddleware(router)

	RegisterRoutes(router)

	server := &http.Server{
		Addr:         cfg.Server.ServerAddr(),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return &App{
		server: server,
		done:   make(chan error),
		cfg:    cfg,
	}
}

func RegisterMiddleware(router *gin.Engine) {

}

func RegisterRoutes(router *gin.Engine) {
	tg := router.Group("/test")
	test.RegisterRoutes(tg)
}

func (a *App) Run() error {

	go a.gracefulShutdown()

	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("listen: %s\n", err)
		return err
	}

	return <-a.done
}
