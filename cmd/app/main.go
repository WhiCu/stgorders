package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/WhiCu/stgorders/internal/app"
	"github.com/WhiCu/stgorders/internal/config"
)

func NewAPP(cfg *config.Config, log *slog.Logger) error {
	app := app.NewApp(cfg, log)

	ctx := context.Background()

	if err := app.Run(ctx); err != nil {
		log.Error("Server failed to start", slog.String("ERR", err.Error()))
		fmt.Printf(`
	==================================
	=                                =
	=    Server failed to start     =
	=                                =
	==================================
%v
		`, err)
		return err
	}
	log.Info("server successfully stopped")
	fmt.Println(`
	==================================
	=                                =
	=    Server successfully stop    =
	=                                =
	==================================
	`)
	return nil
}
