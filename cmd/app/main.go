package main

import (
	"context"
	"log"

	"github.com/WhiCu/stgorders/internal/app"
)

func main() {
	// cfg := config.MustLoadWithEnv()

	app := app.NewApp(nil)

	if err := app.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
	log.Println("Server exited")
}
