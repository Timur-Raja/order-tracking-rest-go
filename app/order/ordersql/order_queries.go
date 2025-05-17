package ordersql

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/timur-raja/order-tracking-rest-go/app/order"
	"github.com/timur-raja/order-tracking-rest-go/db"
)

type GetOrderList struct {
	db.BaseQuery
	Orders []*order.Order
}

func (q *GetOrderList) Exec(ctx context.Context) error {
	if err := pgxscan.Select(ctx, q.DBConn, q.Orders, `
        SELECT id, user_id, created_at, updated_at, deleted_at
        FROM orders
        WHERE deleted_at IS NULL
    `); err != nil {
		return err
	}
	return nil
}

type GetOrderByID struct {
	db.BaseQuery
	Order *order.Order
	Where struct {
		ID int `db:"id"`
	}
}

func (q *GetOrderByID) Exec(ctx context.Context) error {
	if err := pgxscan.Get(ctx, q.DBConn, q.Order, `
        SELECT id, user_id, created_at, updated_at, deleted_at
        FROM orders
        WHERE id = $1
		LIMIT 1
    `); err != nil {
		return err
	}
	return nil
}
