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

type insertOrderQuery struct {
	db.BaseQuery
	redisConn *redis.Client
	Values    struct {
		order.Order
	}
	Returning struct {
		ID int `db:"id"`
	}
}

func NewInsertOrderQuery(conn db.PGExecer, redisConn *redis.Client) *insertOrderQuery {
	return &insertOrderQuery{
		BaseQuery: db.BaseQuery{DBConn: conn},
		redisConn: redisConn,
	}
}

func (q *insertOrderQuery) cacheKey() string {
	return "order:" + strconv.Itoa(q.Values.UserID)
}
func (q *insertOrderQuery) Run(ctx context.Context) error {
	query := `INSERT INTO orders (user_id, status, shipping_address, created_at, updated_at) 
	VALUES($1, $2, $3, $4, $5)
	RETURNING id;`

	if err := pgxscan.Get(ctx, q.DBConn, &q.Returning.ID, query,
		q.Values.UserID,
		q.Values.Status,
		q.Values.ShippingAddress,
		q.Values.CreatedAt,
		q.Values.UpdatedAt); err != nil {
		return err
	}

	// invalidate list of all orders for the user - even if we currently don't have an endpoint to fetch orders by user, this is what would make more sense in a realistic scenario
	if err := q.redisConn.Del(ctx, q.cacheKey()).Err(); err != nil {
		return nil
	}

	return nil
}

// ///////////////////////////////////////////////////////////////////

// no caching since orders could be made frequently and we'd just constantly invalidate the cache
type selectOrderListQuery struct {
	db.BaseQuery
	Orders []*order.Order
}

func NewSelectOrderListQuery(conn db.PGExecer) *selectOrderListQuery {
	return &selectOrderListQuery{
		BaseQuery: db.BaseQuery{DBConn: conn},
		Orders:    []*order.Order{},
	}
}

func (q *selectOrderListQuery) Run(ctx context.Context) error {
	query := `SELECT *
        FROM orders
        WHERE deleted_at IS NULL`

	if err := pgxscan.Select(ctx, q.DBConn, q.Orders, query); err != nil {
		return err
	}
	return nil
}

// ///////////////////////////////////////////////////////////////////
type selectOrderByIDQuery struct {
	db.BaseQuery
	redisConn *redis.Client
	*order.Order
	Where struct {
		ID int `db:"id"`
	}
}

func NewSelectOrderByIDQuery(conn db.PGExecer, redisConn *redis.Client) *selectOrderByIDQuery {
	return &selectOrderByIDQuery{
		BaseQuery: db.BaseQuery{DBConn: conn},
		redisConn: redisConn,
		Order:     &order.Order{},
	}
}

func (q *selectOrderByIDQuery) cacheKey() string {
	return "order:" + strconv.Itoa(q.Where.ID)
}

func (q *selectOrderByIDQuery) Run(ctx context.Context) error {
	if q.redisConn == nil {
		return app.ErrNoRedisConnection.Err
	}
	// check if the order is cached in redis
	// if cached, return the cached order
	if found, err := cache.Get(q.redisConn, ctx, q.cacheKey(), q.Order); err != nil {
		return err
	} else if found {
		return nil
	}

	// not found, run sql query
	query := `SELECT *
        FROM orders
        WHERE id = $1`

	if err := pgxscan.Get(ctx, q.DBConn, q.Order, query,
		q.Where.ID); err != nil {
		return err
	}

	// cache the order in redis
	// order can change status or shipping address, unfrequent operations and we'll invalidate the cache if it happens, so this can be stored for long periods
	if err := cache.Set(q.redisConn, ctx, q.cacheKey(), q.Order, 24*time.Hour); err != nil {
		return nil
	}
	return nil
}

// ///////////////////////////////////////////////////////////////////

// no caching, would need to be invalidated on each order creation, and update
type selectOrderViewListQuery struct {
	db.BaseQuery
	OrdersView []*order.OrderView
}

func NewSelectOrderViewListQuery(conn db.PGExecer) *selectOrderViewListQuery {
	return &selectOrderViewListQuery{
		BaseQuery:  db.BaseQuery{DBConn: conn},
		OrdersView: []*order.OrderView{},
	}
}

func (q *selectOrderViewListQuery) Run(ctx context.Context) error {
	query := `SELECT *
        FROM orders_view
        WHERE deleted_at IS NULL`

	if err := pgxscan.Select(ctx, q.DBConn, q.OrdersView, query); err != nil {
		return err
	}
	return nil
}

// ///////////////////////////////////////////////////////////////////
type SelectOrderViewByIDQuery struct {
	db.BaseQuery
	redisConn *redis.Client
	*order.OrderView
	Where struct {
		ID int `db:"id"`
	}
}

func NewSelectOrderViewByIDQuery(conn db.PGExecer, redisConn *redis.Client) *SelectOrderViewByIDQuery {
	return &SelectOrderViewByIDQuery{
		BaseQuery: db.BaseQuery{DBConn: conn},
		redisConn: redisConn,
		OrderView: &order.OrderView{},
	}
}

func (q *SelectOrderViewByIDQuery) cacheKey() string {
	return "order_view:" + strconv.Itoa(q.Where.ID)
}

func (q *SelectOrderViewByIDQuery) Run(ctx context.Context) error {
	if q.redisConn == nil {
		return app.ErrNoRedisConnection.Err
	}

	// check if the order is cached in redis
	// if cached, return the cached order
	if found, err := cache.Get(q.redisConn, ctx, q.cacheKey(), q.OrderView); err != nil {
		return err
	} else if found {
		return nil
	}

	// not found, run sql query
	query := `SELECT *
        FROM orders_view
        WHERE id = $1`

	if err := pgxscan.Get(ctx, q.DBConn, q.OrderView, query,
		q.Where.ID); err != nil {
		return err
	}

	// cache the order in redis
	// since the embedded orderItems are fetched separately, this can be cached, with invalidation on order update
	if err := cache.Set(q.redisConn, ctx, q.cacheKey(), q.OrderView, 24*time.Hour); err != nil {
		return nil
	}
	return nil
}

// ///////////////////////////////////////////////////////////////////
type updateOrderQuery struct {
	db.BaseQuery
	redisConn *redis.Client
	Values    struct {
		*order.Order
	}
}

func NewUpdateOrderQuery(conn db.PGExecer, redisConn *redis.Client) *updateOrderQuery {
	return &updateOrderQuery{
		BaseQuery: db.BaseQuery{DBConn: conn},
		redisConn: redisConn,
		Values:    struct{ *order.Order }{Order: &order.Order{}},
	}
}

func (q *updateOrderQuery) cacheKeys() []string {
	keys := []string{}
	keys = append(keys, "order:"+strconv.Itoa(q.Values.ID))
	keys = append(keys, "order_view:"+strconv.Itoa(q.Values.ID))
	return keys
}

func (q *updateOrderQuery) Run(ctx context.Context) error {
	if q.redisConn == nil {
		return app.ErrNoRedisConnection.Err
	}
	query := `UPDATE orders
		SET user_id = $1,
			status = $2,
			shipping_address = $3,
			created_at = $4,
			updated_at = $5
		WHERE id = $6`

	_, err := q.DBConn.Exec(ctx, query,
		q.Values.UserID,
		q.Values.Status,
		q.Values.ShippingAddress,
		q.Values.CreatedAt,
		q.Values.UpdatedAt,
		q.Values.ID,
	)

	// invalidate the cache for the order and order view by id
	for _, key := range q.cacheKeys() {
		if err := q.redisConn.Del(ctx, key).Err(); err != nil {
			return nil
		}
	}
	return err
}
