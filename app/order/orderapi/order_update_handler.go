package orderapi

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/olivere/elastic/v7"
	"github.com/timur-raja/order-tracking-rest-go/app"
	"github.com/timur-raja/order-tracking-rest-go/app/order"
	"github.com/timur-raja/order-tracking-rest-go/app/order/orderesrc"
	"github.com/timur-raja/order-tracking-rest-go/app/order/ordersql"
	"github.com/timur-raja/order-tracking-rest-go/app/product/prodsql"
)

//for the order update handler we allow just changing the shipping address
// and cancelling the order

type orderUpdateHandler struct {
	db      *pgxpool.Pool
	es      *elastic.Client
	params  *order.OrderUpdateParams
	orderID int
	Order   *order.Order
}

func OrderUpdateHandler(db *pgxpool.Pool, es *elastic.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &orderUpdateHandler{
			db:     db,
			es:     es,
			params: &order.OrderUpdateParams{},
			Order:  &order.Order{},
		}
		h.Exec(c)
	}
}

func (h *orderUpdateHandler) Exec(c *gin.Context) {
	if err := h.prepare(c); err != nil {
		app.AbortWithErrorResponse(c, app.ErrFailedToLoadParams, err)
		return
	}

	if h.params.Status != nil && *h.params.Status == order.StatusCancelled.String() {
		if h.Order.Status == order.StatusCancelled.String() {
			c.JSON(http.StatusNoContent, nil) // no change
			return
		}
		if err := h.cancelOrder(c); err != nil {
			app.AbortWithErrorResponse(c, app.ErrServerError, err)
			return
		}
	} else if h.params.ShippingAddress != nil {
		// shipping address error handling
		switch *h.params.ShippingAddress {
		case "":
			app.AbortWithErrorResponse(c, order.ErrShippingAddressRequired, nil)
			return
		case h.Order.ShippingAddress:
			c.JSON(http.StatusNoContent, nil) // no change
			return
		default:
			// update the order with the new shipping address
			h.Order.ShippingAddress = *h.params.ShippingAddress
			h.Order.UpdatedAt = time.Now()
			query := ordersql.NewUpdateOrderQuery(h.db)
			query.Values.Order = h.Order
			if err := query.Run(c); err != nil {
				app.AbortWithErrorResponse(c, app.ErrServerError, err)
				return
			}
		}
	}
	// fetch the order view with the updated shipping address
	query := ordersql.NewSelectOrderViewByIDQuery(h.db)
	query.Where.ID = h.orderID
	if err := query.Run(c); err != nil {
		app.AbortWithErrorResponse(c, app.ErrServerError, err)
		return
	}

	// fetch the order items view to build the order
	query2 := ordersql.NewSelectOrderItemViewListByOrderIDQuery(h.db)
	query2.Where.OrderID = h.orderID
	if err := query2.Run(c); err != nil {
		app.AbortWithErrorResponse(c, app.ErrServerError, err)
		return
	}
	query.OrderView.OrderItems = query2.Items
	for _, item := range query2.Items {
		query.OrderView.TotalPrice += item.ItemsPrice
	}

	// index the order in elasticsearch
	indexer := orderesrc.NewOrderIndexer(h.es, "orders")
	if err := indexer.Run(c, query.OrderView); err != nil {
		app.AbortWithErrorResponse(c, app.ErrServerError, err)
		return
	}
	c.JSON(http.StatusOK, query.OrderView)
}

func (h *orderUpdateHandler) cancelOrder(c *gin.Context) error {
	// fetch related order items
	query := ordersql.NewSelectOrderItemListByOrderIDQuery(h.db)
	query.Where.OrderID = h.orderID
	if err := query.Run(c); err != nil {
		app.AbortWithErrorResponse(c, app.ErrServerError, err)
		return app.ErrServerError.Err
	}

	var prodIDs []int
	prodQuantityMap := make(map[int]int)
	for _, item := range query.Items {
		prodIDs = append(prodIDs, item.ProductID)
		prodQuantityMap[item.ProductID] = item.Quantity
	}

	// start a transaction to update the order status and return the stock to the products
	tx, err := h.db.BeginTx(c, pgx.TxOptions{})
	if err != nil {
		app.AbortWithErrorResponse(c, app.ErrServerError, err)
		return app.ErrServerError.Err
	}

	h.Order.Status = order.StatusCancelled.String()
	h.Order.UpdatedAt = time.Now()

	// update the order status
	query2 := ordersql.NewUpdateOrderQuery(tx)
	query2.Values.Order = h.Order
	if err := query2.Run(c); err != nil {
		tx.Rollback(c)
		app.AbortWithErrorResponse(c, app.ErrServerError, err)
		return app.ErrServerError.Err
	}

	// fetch the product list and lock for update to avoid concurrency issues
	query3 := prodsql.NewSelectProductListByIDsForUpdateQuery(tx)
	query3.Where.IDs = prodIDs
	if err := query3.Run(c); err != nil {
		tx.Rollback(c)
		app.AbortWithErrorResponse(c, app.ErrServerError, err)
		return app.ErrServerError.Err
	}

	// return the stock to the products
	for _, item := range query3.Products {
		item.Stock += prodQuantityMap[item.ID]
		query4 := prodsql.NewUpdateProductStockQuery(tx)
		query4.Values.Stock = item.Stock
		query4.Where.ID = item.ID
		if err := query4.Run(c); err != nil {
			tx.Rollback(c)
			app.AbortWithErrorResponse(c, app.ErrServerError, err)
			return app.ErrServerError.Err
		}
	}

	if err := tx.Commit(c); err != nil {
		app.AbortWithErrorResponse(c, app.ErrServerError, err)
		return app.ErrServerError.Err
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// utils

// load params and the requested order from the database
func (h *orderUpdateHandler) prepare(c *gin.Context) error {
	idParam := c.Param("order_id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		app.AbortWithErrorResponse(c, app.ErrFailedToLoadParams, err)
		return app.ErrFailedToLoadParams.Err
	}

	h.orderID = id
	if err := c.ShouldBindJSON(h.params); err != nil {
		app.AbortWithErrorResponse(c, app.ErrFailedToLoadParams, err)
		return app.ErrFailedToLoadParams.Err
	}

	// load order if it exists
	query := ordersql.NewSelectOrderByIDQuery(h.db)
	query.Where.ID = h.orderID
	if err := query.Run(c); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			app.AbortWithErrorResponse(c, order.ErrOrderNotFound, err)
			return order.ErrOrderNotFound.Err
		} else {
			app.AbortWithErrorResponse(c, app.ErrServerError, err)
			return app.ErrServerError.Err
		}
	}

	h.Order = query.Order
	return nil
}
