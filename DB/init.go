package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/timur-raja/order-tracking-rest-go/app"
)

func Init(cfg *app.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)
	return pgxpool.Connect(context.Background(), dsn)
}

type Query interface { // for mocking
	Exec(ctx context.Context) error
}

type BaseQuery struct {
	DBConn *pgxpool.Pool
}
