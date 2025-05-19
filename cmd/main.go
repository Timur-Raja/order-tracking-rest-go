package main

import (
	"fmt"
	"log"

	"github.com/timur-raja/order-tracking-rest-go/api"
	"github.com/timur-raja/order-tracking-rest-go/config"
	"github.com/timur-raja/order-tracking-rest-go/db"

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

	dbConn, err := db.Init(cfg.DB.DSN)
	if err != nil {
		return err
	}
	defer dbConn.Close()

	server := gin.New()
	server.Use(gin.Recovery(), api.ErrorLogger())

	r := api.NewRouter(dbConn)
	r.Setup(server)

	addr := fmt.Sprintf("%s:%s", cfg.WebServer.Host, cfg.WebServer.Port)
	return server.Run(addr)
}
