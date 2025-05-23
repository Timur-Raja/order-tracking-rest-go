package ordersql

import (
	"context"
	"strconv"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/go-redis/redis/v8"
	"github.com/timur-raja/order-tracking-rest-go/app"
	"github.com/timur-raja/order-tracking-rest-go/app/order"
	"github.com/timur-raja/order-tracking-rest-go/cache"
	"github.com/timur-raja/order-tracking-rest-go/db"
)

type insertOrderItemQuery struct { // never used, just added as general code showcase
	db.BaseQuery
	Values struct {
		order.OrderItem
	}
	Returning struct {
		ID int `db:"id"`
	}
}

func NewInsertOrderItemQuery(conn db.PGExecer) *insertOrderItemQuery {
	return &insertOrderItemQuery{
		BaseQuery: db.BaseQuery{DBConn: conn},
	}
}

func (q *insertOrderItemQuery) Run(ctx context.Context) error {
	query := `INSERT INTO order_items (order_id, product_id, quantity, price)
	VALUES($1, $2, $3)
	RETURNING token`

	if err := pgxscan.Get(ctx, q.DBConn, &q.Returning.ID, query,
		q.Values.OrderID,
		q.Values.ProductID,
		q.Values.Quantity); err != nil {
		return err
	}
	return nil
}

// ////////////////////////////////////////////////////////////////////////////

// no caching invalidation, we do not have a generic select list of order items which would need to be invalidated
type insertOrderItemsListQuery struct {
	db.BaseQuery
	Values struct {
		Items []*order.OrderItem
	}
}

func NewInsertOrderItemsListQuery(conn db.PGExecer) *insertOrderItemsListQuery {
	return &insertOrderItemsListQuery{
		BaseQuery: db.BaseQuery{DBConn: conn},
	}
}

func (q *insertOrderItemsListQuery) Run(ctx context.Context) error {
	if len(q.Values.Items) == 0 {
		return nil
	}

	query := `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES `

	rows := q.Values.Items
	fields := []string{"OrderID", "ProductID", "Quantity", "Price"}

	// func to dynamically create value strings
	valueStrings, args := db.CreateValueStrings(rows, fields)
	query = query + valueStrings

	_, err := q.DBConn.Exec(ctx, query, args...)
	return err
}

// ///////////////////////////////////////////////////////////////////
type selectOrderItemListByOrderIDQuery struct {
	db.BaseQuery
	redisConn *redis.Client
	Where     struct {
		OrderID int `db:"order_id"`
	}
	Items []*order.OrderItem
}

func NewSelectOrderItemListByOrderIDQuery(dbConn db.PGExecer, redisConn *redis.Client) *selectOrderItemListByOrderIDQuery {
	return &selectOrderItemListByOrderIDQuery{
		BaseQuery: db.BaseQuery{DBConn: dbConn},
		redisConn: redisConn,
		Items:     []*order.OrderItem{},
	}
}

func (q *selectOrderItemListByOrderIDQuery) cacheKey() string {
	return "order_items:order:" + strconv.Itoa(q.Where.OrderID)
}

func (q *selectOrderItemListByOrderIDQuery) Run(ctx context.Context) error {
	if q.redisConn == nil {
		return app.ErrNoRedisConnection.Err
	}

	// check if the order items are cached in redis
	// if cached, return the cached items
	if found, err := cache.Get(q.redisConn, ctx, q.cacheKey(), &q.Items); err != nil {
		return err
	} else if found {
		return nil
	}

	// not found, rund sql query
	query := `SELECT *
        FROM order_items
        WHERE order_id = $1`

	if err := pgxscan.Select(ctx, q.DBConn, &q.Items, query,
		q.Where.OrderID,
	); err != nil {
		return err
	}

	// cache the order items in redis
	if err := cache.Set(q.redisConn, ctx, q.cacheKey(), q.Items, 24*time.Hour); err != nil { // order items shouldn't change, so we can cache for a long period
		return err
	}
	return nil
}

// ////////////////////////////////////////////////////////////////////

// no caching, the view is depdendent on real time product stock changes, which might update very often and potentially reference thousands of records. Making ccaching very inefficient
type selectOrderItemViewListByOrderIDQuery struct {
	db.BaseQuery
	Where struct {
		OrderID int `db:"order_id"`
	}
	Items []*order.OrderItemView
}

func NewSelectOrderItemViewListByOrderIDQuery(conn db.PGExecer) *selectOrderItemViewListByOrderIDQuery {
	return &selectOrderItemViewListByOrderIDQuery{
		BaseQuery: db.BaseQuery{DBConn: conn},
		Items:     []*order.OrderItemView{},
	}
}

func (q *selectOrderItemViewListByOrderIDQuery) Run(ctx context.Context) error {
	query := `
		SELECT *
		FROM order_items_view
		WHERE order_id = $1
	`
	if err := pgxscan.Select(ctx, q.DBConn, &q.Items, query, q.Where.OrderID); err != nil {
		return err
	}
	return nil
}
