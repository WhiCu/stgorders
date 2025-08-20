package main

import (
	"context"
	"fmt"
	"log"

	"github.com/WhiCu/stgorders/internal/app"
	"github.com/WhiCu/stgorders/internal/config"
)

func main() {
	cfg := config.MustLoadWithDefault("./config/config.yaml")

	app := app.NewApp(cfg)

	ctx := context.Background()

	if err := app.Run(ctx); err != nil {
		log.Fatalf(`
	==================================
	=                                =
	=    Server failed to start     =
	=                                =
	==================================
%v
		`, err)
	}
	fmt.Println(`
	==================================
	=                                =
	=    Server successfully stop    =
	=                                =
	==================================
	`)
}
