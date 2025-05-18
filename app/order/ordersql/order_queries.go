package ordersql

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/timur-raja/order-tracking-rest-go/app/order"
	"github.com/timur-raja/order-tracking-rest-go/db"
)

type insertOrderQuery struct {
	db.BaseQuery
	Values struct {
		order.Order
	}
	Returning struct {
		ID int `db:"id"`
	}
}

func NewInsertOrderQuery(conn db.PGExecer) *insertOrderQuery {
	return &insertOrderQuery{
		BaseQuery: db.BaseQuery{DBConn: conn},
	}
}

func (q *insertOrderQuery) Run(ctx context.Context) error {
	query := `INSERT INTO orders (user_id, status, created_at, updated_at) 
	VALUES($1, $2, $3, $4)
	RETURNING token;`

	if err := pgxscan.Get(ctx, q.DBConn, &q.Returning.ID, query,
		q.Values.UserID,
		q.Values.Status,
		q.Values.CreatedAt,
		q.Values.UpdatedAt); err != nil {
		return err
	}
	return nil
}

// ///////////////////////////////////////////////////////////////////
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
	*order.Order
	Where struct {
		ID int `db:"id"`
	}
}

func NewSelectOrderByIDQuery(conn db.PGExecer) *selectOrderByIDQuery {
	return &selectOrderByIDQuery{
		BaseQuery: db.BaseQuery{DBConn: conn},
	}
}

func (q *selectOrderByIDQuery) Run(ctx context.Context) error {
	query := `SELECT *
        FROM orders
        WHERE id = $1`

	if err := pgxscan.Select(ctx, q.DBConn, q.Order, query,
		q.Where.ID); err != nil {
		return err
	}
	return nil
}

// ///////////////////////////////////////////////////////////////////
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
	*order.OrderView
	Where struct {
		ID int `db:"id"`
	}
}

func NewSelectOrderViewByIDQuery(conn db.PGExecer) *SelectOrderViewByIDQuery {
	return &SelectOrderViewByIDQuery{
		BaseQuery: db.BaseQuery{DBConn: conn},
	}
}

func (q *SelectOrderViewByIDQuery) Run(ctx context.Context) error {
	query := `SELECT *
        FROM orders
        WHERE id = $1`

	if err := pgxscan.Select(ctx, q.DBConn, q.OrderView, query,
		q.Where.ID); err != nil {
		return err
	}
	return nil
}
