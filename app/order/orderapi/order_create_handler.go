package orderapi

import (
	"time"

	"github.com/timur-raja/order-tracking-rest-go/app"
	"github.com/timur-raja/order-tracking-rest-go/app/order"
	"github.com/timur-raja/order-tracking-rest-go/app/order/orderesrc"
	"github.com/timur-raja/order-tracking-rest-go/app/order/ordersql"
	"github.com/timur-raja/order-tracking-rest-go/app/product"
	"github.com/timur-raja/order-tracking-rest-go/app/product/prodsql"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/olivere/elastic/v7"
)

type orderCreateHandler struct {
	db *pgxpool.Pool
	es *elastic.Client

	// req
	userID int
	params *order.OrderCreateParams

	// helpers
	itemsMap map[int]int
	prodIDs  []int
	products []*product.Product

	// insert
	items []*order.OrderItem
}

func OrderCreateHandler(db *pgxpool.Pool, es *elastic.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &orderCreateHandler{
			db:       db,
			es:       es,
			params:   &order.OrderCreateParams{},
			itemsMap: make(map[int]int),
			prodIDs:  []int{},
			products: []*product.Product{},
			items:    []*order.OrderItem{},
		}
		h.Exec(c)
	}
}

func (h *orderCreateHandler) Exec(c *gin.Context) {
	// setup the struct
	if err := h.Prepare(c); err != nil {
		app.AbortWithErrorResponse(c, app.ErrFailedToLoadParams, err)
		return
	}

	// load the product list from the database and lock for update to avoid cocurrency issues
	// this is needed because we need to update the product stock and avoid inconsistent data

	// start a transaction
	tx, err := h.db.BeginTx(c, pgx.TxOptions{})
	if err != nil {
		app.AbortWithErrorResponse(c, app.ErrServerError, err)
		return
	}
	query := prodsql.NewSelectProductListByIDsForUpdateQuery(h.db)
	query.Where.IDs = h.prodIDs
	if err := query.Run(c); err != nil {
		tx.Rollback(c)
		app.AbortWithErrorResponse(c, app.ErrServerError, err)
		return
	}

	// check if all products exist
	if len(query.Products) != len(h.prodIDs) {
		tx.Rollback(c)
		app.AbortWithErrorResponse(c, product.ErrProductsNotFound, product.ErrProductsNotFound.Err)
		return
	}

	h.products = query.Products

	// check if all products are valid
	if err := h.ValidateAndCreateItems(); err != nil {
		tx.Rollback(c)
		app.AbortWithErrorResponse(c, product.ErrProductsNotFound, err)
		return
	}

	// insert order items and orders as part of the transaction
	query2 := ordersql.NewInsertOrderQuery(tx)
	query2.Values.UserID = h.userID
	query2.Values.Status = order.StatusCreated.String()
	query2.Values.ShippingAddress = h.params.ShippingAddress
	query2.Values.CreatedAt = time.Now()
	query2.Values.UpdatedAt = time.Now()

	if err := query2.Run(c); err != nil {
		tx.Rollback(c)
		app.AbortWithErrorResponse(c, app.ErrServerError, err)
		return
	}

	orderID := query2.Returning.ID

	// return stock to products
	productStock := prodsql.ProductStock{}
	productStockList := make([]prodsql.ProductStock, len(h.products))

	for _, item := range h.products {
		productStock.ID = item.ID
		productStock.Stock = item.Stock + h.itemsMap[item.ID]
		productStockList = append(productStockList, productStock)
	}

	query3 := prodsql.NewUpdateProductsStockByIDsQuery(tx)
	query3.Values.ProductStockList = productStockList
	if err := query3.Run(c); err != nil {
		tx.Rollback(c)
		app.AbortWithErrorResponse(c, app.ErrServerError, err)
		return
	}

	// set the order id for each item
	for _, item := range h.items {
		item.OrderID = orderID
	}

	// bulk insert order items
	query4 := ordersql.NewInsertOrderItemsListQuery(tx)
	query4.Values.Items = h.items
	if err := query4.Run(c); err != nil {
		tx.Rollback(c)
		app.AbortWithErrorResponse(c, app.ErrServerError, err)
		return
	}
	if err := tx.Commit(c); err != nil {
		app.AbortWithErrorResponse(c, app.ErrServerError, err)
		return
	}

	// return the order
	query5 := ordersql.NewSelectOrderViewByIDQuery(h.db)
	query5.Where.ID = orderID
	if err := query5.Run(c); err != nil {
		app.AbortWithErrorResponse(c, app.ErrServerError, err)
		return
	}

	// fetch the items related to the order
	query6 := ordersql.NewSelectOrderItemViewListByOrderIDQuery(h.db)
	query6.Where.OrderID = orderID
	if err := query6.Run(c); err != nil {
		app.AbortWithErrorResponse(c, app.ErrServerError, err)
		return
	}

	query5.OrderView.OrderItems = query6.Items
	for _, item := range query6.Items {
		query5.OrderView.TotalPrice += item.ItemsPrice
	}

	// index order in elasticsearch
	indexer := orderesrc.NewOrderIndexer(h.es, "orders")
	if err := indexer.Run(c, query5.OrderView); err != nil {
		app.AbortWithErrorResponse(c, app.ErrServerError, err)
		return
	}

	c.JSON(201, query5.OrderView)
}

// load info and initialize helpers
func (h *orderCreateHandler) Prepare(c *gin.Context) error {
	// get the user id from the context
	id, ok := c.Get("user_id")
	if !ok {
		return app.ErrAuthenticationRequired.Err
	}
	h.userID = id.(int)

	// load the order params from the request body
	if err := c.ShouldBindJSON(h.params); err != nil {
		return err
	}

	if len(h.params.OrderItems) == 0 {
		return order.ErrOrderItemsRequired.Err
	}

	h.params.Sanitize()

	for _, orderItem := range h.params.OrderItems {
		h.itemsMap[orderItem.ProductID] = orderItem.Quantity
		h.prodIDs = append(h.prodIDs, orderItem.ProductID)
	}
	return nil
}

func (h *orderCreateHandler) ValidateAndCreateItems() error {
	for _, item := range h.products {
		if quantity, ok := h.itemsMap[item.ID]; !ok {
			return product.ErrProductsNotFound.Err
		} else {
			orderItem := &order.OrderItem{
				ProductID: item.ID,
				Quantity:  quantity,
				Price:     item.Price * float32(quantity),
			}
			h.items = append(h.items, orderItem)
		}
	}
	return nil
}
