package testing

import (
	"context"
	"database/sql"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/timur-raja/order-tracking-rest-go/api"
	"github.com/timur-raja/order-tracking-rest-go/config"
	"github.com/timur-raja/order-tracking-rest-go/db"
)

//TODO: needs imprpovements, very rudimentary implementation due to time constraints

var baseURL string

func TestMain(m *testing.M) {
	// load configs
	if err := os.Chdir(".."); err != nil {
		log.Fatalf("could not chdir to project root: %v", err)
	}
	cfg := new(config.Config)
	if err := cfg.LoadConfig(); err != nil {
		log.Fatalf("loading envs: %v", err)
	}

	// run migrations on test DB
	sqlDB, err := sql.Open("pgx", cfg.TestDB.DSN)
	if err != nil {
		log.Fatalf("opening migration DB: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("pinging migration DB: %v", err)
	}
	if err := db.MigrateUp(sqlDB, "db/migrations"); err != nil {
		log.Fatalf("running migrations: %v", err)
	}

	// seed user and session
	if _, err := sqlDB.Exec(`
 	 INSERT INTO users (name, email, password) VALUES
    ('John Doe','john.doe@example.com','hashed_password_here');
 	 INSERT INTO user_sessions(token, user_id) VALUES
    ('test123', (SELECT id FROM users WHERE email='john.doe@example.com'));
`); err != nil {
		log.Fatalf("seeding test user/session failed: %v", err)
	}

	sqlDB.Close()

	dbConn, err := db.Init(cfg.TestDB.DSN)
	if err != nil {
		log.Fatalf("connecting to test db: %v", err)
	}
	defer dbConn.Close()

	// setup test gin server
	server := gin.New()
	server.Use(gin.Recovery(), api.ErrorLogger())

	r := api.NewRouter(dbConn)
	r.Setup(server)

	ts := httptest.NewServer(server)
	baseURL = ts.URL
	defer ts.Close()

	// run tests
	code := m.Run()

	ctx := context.Background()
	// clean affected table for subequent tests
	if _, err := dbConn.Exec(ctx, `
       TRUNCATE order_items, orders, user_sessions, users
      RESTART IDENTITY CASCADE
   `); err != nil {
		log.Fatalf("cleanup – truncate failed: %v", err)
	}

	os.Exit(code)
}
