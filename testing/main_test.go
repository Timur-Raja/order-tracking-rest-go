package testing

import (
	"database/sql"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/timur-raja/order-tracking-rest-go/api"
	"github.com/timur-raja/order-tracking-rest-go/config"
	"github.com/timur-raja/order-tracking-rest-go/db"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/olivere/elastic/v7"
)

var baseURL string

func TestMain(m *testing.M) {
	// point tests at the test‐only ES on 9201
	os.Setenv("ES_URL", "http://localhost:9201")

	// load config (now cfg.ES.URL == "http://localhost:9201")
	if err := os.Chdir(".."); err != nil {
		log.Fatalf("cd: %v", err)
	}
	cfg := new(config.Config)
	if err := cfg.LoadConfig(); err != nil {
		log.Fatalf("load config: %v", err)
	}

	// migrate & connect Postgres
	sqlDB, err := sql.Open("pgx", cfg.TestDB.DSN)
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

	pgPool, err := db.Init(cfg.TestDB.DSN)
	if err != nil {
		log.Fatalf("init pg pool: %v", err)
	}
	defer pgPool.Close()

	// init ES test service
	esClient, err := elastic.NewClient(
		elastic.SetURL(cfg.TestES.URL),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)
	if err != nil {
		log.Fatalf("init ES client: %v", err)
	}

	// start up your Gin router
	gin.SetMode(gin.TestMode)
	server := gin.New()
	server.Use(gin.Recovery(), api.ErrorLogger())

	// pass both pg and es into your router
	r := api.NewRouter(pgPool, esClient)
	r.Setup(server)

	ts := httptest.NewServer(server)
	baseURL = ts.URL
	defer ts.Close()

	code := m.Run()
	//clean testing db for next run
	if err := db.MigrateDrop(sqlDB, "db/migrations"); err != nil {
		log.Fatalf("cleaning dn migrations: %v", err)
	}
	sqlDB.Close()
	os.Exit(code)
}
