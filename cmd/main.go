package main

import (
	"log"

	"github.com/timur-raja/order-tracking-rest-go/api"
	"github.com/timur-raja/order-tracking-rest-go/app"
	"github.com/timur-raja/order-tracking-rest-go/config"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Fatal startup error: %v", err)
	}
}

func run() error {
	cfg := new(config.Config)
	if err := cfg.LoadConfig(); err != nil {
		return err
	}

	services := app.InitServices(*cfg, false)

	server := gin.New()
	server.Use(gin.Recovery(), api.ErrorLogger())

	r := api.NewRouter(services)
	r.Setup(server)

	addr := cfg.WebServer.URL
	return server.Run(addr)
}
