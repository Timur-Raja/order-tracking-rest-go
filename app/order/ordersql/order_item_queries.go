package ordersql

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/timur-raja/order-tracking-rest-go/app/order"
	"github.com/timur-raja/order-tracking-rest-go/db"
)

type insertOrderItemQuery struct {
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

	_, err := q.BaseQuery.DBConn.Exec(ctx, query, args...)
	return err
}

// ///////////////////////////////////////////////////////////////////
type selectOrderItemListByOrderIDQuery struct {
	db.BaseQuery
	Where struct {
		OrderID   int `db:"order_id"`
		ProductID int `db:"product_id"`
	}
	Orders []*order.Order
}

func NewSelectOrderItemListByOrderIDQuery(conn db.PGExecer) *selectOrderItemListByOrderIDQuery {
	return &selectOrderItemListByOrderIDQuery{
		BaseQuery: db.BaseQuery{DBConn: conn},
		Orders:    []*order.Order{},
	}
}

func (q *selectOrderItemListByOrderIDQuery) Run(ctx context.Context) error {
	query := `SELECT *
        FROM order_items
        WHERE order_id = $1
		AND product_id = $2`

	if err := pgxscan.Select(ctx, q.DBConn, q.Orders, query,
		q.Where.OrderID,
		q.Where.ProductID,
	); err != nil {
		return err
	}
	return nil
}

// ///////////////////////////////////////////////////////////////////////
type updateOrderItemQuery struct {
	db.BaseQuery
	Values struct {
		order.OrderItem
	}
}

func NewUpdateOrderItemQuery(conn db.PGExecer) *updateOrderItemQuery {
	return &updateOrderItemQuery{
		BaseQuery: db.BaseQuery{DBConn: conn},
	}
}

func (q *updateOrderItemQuery) Run(ctx context.Context) error {
	query := `
      UPDATE order_items
         SET order_id   = $1,
             product_id = $2,
             quantity   = $3,
             price      = $4
       WHERE order_id = $5
	   AND product_id = $6
    `
	if _, err := q.DBConn.Exec(ctx, query,
		q.Values.OrderID,
		q.Values.ProductID,
		q.Values.Quantity,
		q.Values.Price,
		q.Values.OrderID,
		q.Values.ProductID,
	); err != nil {
		return err
	}
	return nil
}

// ////////////////////////////////////////////////////////////////////
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
