package app

import (
	"context"
	"fmt"

	"github.com/timur-raja/order-tracking-rest-go/config"
	"github.com/timur-raja/order-tracking-rest-go/db"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/olivere/elastic/v7"
)

type Services struct {
	DB    *pgxpool.Pool
	ES    *elastic.Client
	Redis *redis.Client
}

func InitServices(cfg config.Config, isTesting bool) *Services {
	// Initialize envs based on environment
	dbDSN := cfg.DB.DSN
	redisURL := cfg.Redis.URL
	esURL := cfg.ES.URL
	if isTesting {
		dbDSN = cfg.DB.TestDSN
		redisURL = cfg.Redis.TestURL
		esURL = cfg.ES.TestURL
	}

	// postgres
	dbConn, err := db.Init(dbDSN)
	if err != nil {
		panic(err)
	}
	if err := dbConn.Ping(context.Background()); err != nil {
		panic(fmt.Sprintf("Cannot connect to Postgres: %v", err))
	}

	// elasticsearch
	esConn, err := elastic.NewClient(
		elastic.SetURL(esURL),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)

	if err != nil {
		panic(err)
	}
	_, _, err = esConn.Ping(esURL).Do(context.Background())
	if err != nil {
		panic(fmt.Sprintf("Cannot connect to Elasticsearch: %v", err))
	}

	// redis
	redisConn := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	if err := redisConn.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Sprintf("Cannot connect to Redis: %v", err))
	}
	return &Services{
		DB:    dbConn,
		ES:    esConn,
		Redis: redisConn,
	}
}
