package testing

import (
	"context"
	"database/sql"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/timur-raja/order-tracking-rest-go/api"
	"github.com/timur-raja/order-tracking-rest-go/app"
	"github.com/timur-raja/order-tracking-rest-go/config"
	"github.com/timur-raja/order-tracking-rest-go/db"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var baseURL string

func TestMain(m *testing.M) {
	// point tests at the test‐only ES on 9201
	os.Setenv("ES_URL", "http://localhost:9201")

	// load config
	if err := os.Chdir(".."); err != nil {
		log.Fatalf("cd: %v", err)
	}
	cfg := new(config.Config)
	if err := cfg.LoadConfig(); err != nil {
		log.Fatalf("load config: %v", err)
	}

	// run migrations on testdb
	sqlDB, err := sql.Open("pgx", cfg.DB.TestDSN)
	if err != nil {
		log.Fatalf("open test db: %v", err)
	}
	if err := db.MigrateUp(sqlDB, "db/migrations"); err != nil {
		log.Fatalf("migrate up: %v", err)
	}
	defer func() {
		db.MigrateDrop(sqlDB, "db/migrations")
		sqlDB.Close()
	}()

	// setup test services
	services := app.InitServices(*cfg, true)

	// start up your Gin router
	gin.SetMode(gin.TestMode)
	server := gin.New()
	server.Use(gin.Recovery(), api.ErrorLogger())

	// pass both pg and es into your router
	r := api.NewRouter(services)
	r.Setup(server)

	ts := httptest.NewServer(server)
	baseURL = ts.URL
	defer ts.Close()

	code := m.Run()
	//clean testing dbs for next run
	if err := db.MigrateDrop(sqlDB, "db/migrations"); err != nil {
		log.Fatalf("cleaning dn migrations: %v", err)
	}
	sqlDB.Close()

	if err := services.Redis.FlushAll(context.Background()); err != nil {
		log.Printf("flushing redis: %v", err)
	}

	if _, err := services.ES.DeleteIndex("orders").Do(context.Background()); err != nil {
		log.Fatalf("deleting elasticsearch indexes: %v", err)
	}
	os.Exit(code)
}
