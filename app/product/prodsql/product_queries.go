package prodsql

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/timur-raja/order-tracking-rest-go/app/product"
	"github.com/timur-raja/order-tracking-rest-go/db"
)

type selectProductByIDQuery struct {
	db.BaseQuery
	Where struct {
		ID int
	}
	*product.Product
}

func NewSelectProductByIDQuery(conn db.PGExecer) *selectProductByIDQuery {
	return &selectProductByIDQuery{
		BaseQuery: db.BaseQuery{DBConn: conn},
		Product:   &product.Product{},
	}
}

func (q *selectProductByIDQuery) Run(ctx context.Context) error {
	query := `
	SELECT *
	FROM products
	WHERE id = $1`

	err := pgxscan.Get(ctx, q.DBConn, q.Product, query,
		q.Where.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

// ///////////////////////////////////////////////////////////////////
// selects rows and locks them for update, to prevent concurrent updates
// only used when we need to fetch products to then update them (when an order is placed)
type selectProductListByIDsForUpdateQuery struct {
	db.BaseQuery
	Where struct {
		IDs []int
	}
	Products []*product.Product
}

func NewSelectProductListByIDsForUpdateQuery(conn db.PGExecer) *selectProductListByIDsForUpdateQuery {
	return &selectProductListByIDsForUpdateQuery{
		BaseQuery: db.BaseQuery{DBConn: conn},
		Products:  []*product.Product{},
	}
}

func (q *selectProductListByIDsForUpdateQuery) Run(ctx context.Context) error {
	query := `
	SELECT *
	FROM products
	WHERE id = ANY($1)
	FOR UPDATE`

	if err := pgxscan.Select(ctx, q.DBConn, &q.Products, query, q.Where.IDs); err != nil {
		return err
	}
	return nil
}

// / ///////////////////////////////////////////////////////////////////
type UpdateProductStockByIDQuery struct {
	db.BaseQuery
	Values struct {
		Stock int
	}
	Where struct {
		ID int
	}
}

func NewUpdateProductStockQuery(conn db.PGExecer) *UpdateProductStockByIDQuery {
	return &UpdateProductStockByIDQuery{
		BaseQuery: db.BaseQuery{DBConn: conn},
	}
}

func (q *UpdateProductStockByIDQuery) Run(ctx context.Context) error {
	query := `
        UPDATE products
        SET stock = $1
        WHERE id = $2
    `
	_, err := q.DBConn.Exec(ctx, query, q.Values.Stock, q.Where.ID)
	return err
}
