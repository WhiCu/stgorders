package main

import (
	"fmt"
	"log"

	"github.com/WhiCu/stgorders/internal/app"
	"github.com/WhiCu/stgorders/internal/config"
)

func main() {
	cfg := config.MustLoadWithEnv()

	app := app.NewApp(cfg)

	fmt.Printf(`
=============================================

Server is running on %s

=============================================
Configuration:

%s

=============================================
`, cfg.Server.ServerAddr(), cfg.Format())

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
	log.Println("Server exited")
}
