package db

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/timur-raja/order-tracking-rest-go/config"
)

func Init(cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := cfg.DB.DSN
	return pgxpool.Connect(context.Background(), dsn)
}

type PGExecer interface {
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}
type Query interface { // for mocking
	Exec(ctx context.Context) error
}

type BaseQuery struct {
	DBConn PGExecer
}
