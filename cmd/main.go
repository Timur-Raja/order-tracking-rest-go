package main

import (
	"fmt"
	"log"

	db "github.com/timur-raja/order-tracking-rest-go/DB"
	"github.com/timur-raja/order-tracking-rest-go/app"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Fatal startup error: %v", err)
	}
}

func run() error {
	cfg := new(app.Config)
	if err := cfg.LoadConfig(); err != nil {
		return err
	}

	dbConn, err := db.Init(cfg)
	if err != nil {
		return err
	}
	defer dbConn.Close()

	router := gin.Default()

	r := app.NewRouter(dbConn) // 👈 inject db here
	r.Setup(router)

	addr := fmt.Sprintf("%s:%s", cfg.WebServer.Host, cfg.WebServer.Port)
	return router.Run(addr)
}
