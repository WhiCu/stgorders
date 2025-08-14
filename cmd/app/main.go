package main

import (
	"context"
	"log"

	"github.com/WhiCu/stgorders/internal/app"
	"github.com/WhiCu/stgorders/internal/config"
)

func main() {
	cfg := config.MustLoadWithDefault("./config/config.yaml")

	app := app.NewApp(cfg)

	if err := app.Run(context.Background()); err != nil {
		panic(err)
	}
	log.Println("Server exited")
}
