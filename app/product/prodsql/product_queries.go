package prodsql

import (
	"context"
	"fmt"
	"strings"

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

// ///////////////////////////////////////////////////////////////////
// dynamically generated query to update multiple products at once

type ProductStock struct {
	ID    int
	Stock int
}
type UpdateProductsStockByIDsQuery struct {
	db.BaseQuery
	Values struct {
		ProductStockList []ProductStock
	}
}

func NewUpdateProductsStockByIDsQuery(conn db.PGExecer) *UpdateProductsStockByIDsQuery {
	return &UpdateProductsStockByIDsQuery{
		BaseQuery: db.BaseQuery{DBConn: conn},
	}
}

func (q *UpdateProductsStockByIDsQuery) Run(ctx context.Context) error {
	numberOfUpdates := len(q.Values.ProductStockList)
	if numberOfUpdates == 0 {
		return nil
	}

	// prepapre slices that will hold the dyanamic pieces of the query based on the number of products to update
	caseClauses := make([]string, 0, numberOfUpdates) // holds the WHEN - THEN clauses
	idParams := make([]string, 0, numberOfUpdates)    // holds params to append into the IN clause

	// slices to holld actual values that will be passed to the query in sequential order
	// it is of type any (interface{}) so it can be passed as a series of dynamic length variadic arguments to the Exec func
	args := make([]any, 0, numberOfUpdates*3) // 3 values per product: 2 for the CASE clause and 1 for the IN clause

	// build WHEN - THEN clauses
	for i, v := range q.Values.ProductStockList {
		// set parameter positions: start from 0 and 1 and increment each by 2
		idPos := i*2 + 1
		stockPos := i*2 + 2

		// populate the slice of WHEN - THEN clauses
		caseClauses = append(caseClauses,
			fmt.Sprintf("WHEN $%d THEN $%d", idPos, stockPos),
		)

		// collect the actual values to pass to each clause in sequential order
		args = append(args, v.ID, v.Stock)
	}

	// build list of params to pass to IN clause to indicate all the products to be updated
	offset := numberOfUpdates * 2 // offset to start from the last param number used in the CASE clause
	for _, v := range q.Values.ProductStockList {
		offset++                                                // increment offset for each product
		idParams = append(idParams, fmt.Sprintf("$%d", offset)) // populate the slice of params to append into the IN clause
		args = append(args, v.ID)                               // collect the actual product ID that will be passed to the IN clause
	}

	// build full query
	query := fmt.Sprintf(`
	UPDATE products
	SET stock = CASE id
		%s
		ELSE stock
	END
	WHERE id IN (%s)`,
		strings.Join(caseClauses, "\n    "),
		strings.Join(idParams, ", "),
	)

	_, err := q.DBConn.Exec(ctx, query, args...) // exec with all the values
	return err
}
