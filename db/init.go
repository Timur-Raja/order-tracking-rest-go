package db

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// function to easily intialize db connection
func Init(DSN string) (*pgxpool.Pool, error) {
	return pgxpool.Connect(context.Background(), DSN)
}

// interface so we can pass both pools and transactions to queries to execute the sql statements
// wrapper around pgxpool.Pool and pgx.Tx
type PGExecer interface {
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type BaseQuery struct {
	DBConn PGExecer
}
