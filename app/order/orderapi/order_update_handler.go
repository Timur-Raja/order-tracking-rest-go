package orderapi

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/timur-raja/order-tracking-rest-go/app"
	"github.com/timur-raja/order-tracking-rest-go/app/order"
	"github.com/timur-raja/order-tracking-rest-go/app/order/orderesrc"
	"github.com/timur-raja/order-tracking-rest-go/app/order/ordersql"
	"github.com/timur-raja/order-tracking-rest-go/app/product/prodsql"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

//for the order update handler we allow just changing the shipping address
// and cancelling the order

type orderUpdateHandler struct {
	connections *app.Services
	params      *order.OrderUpdateParams
	orderID     int
	order       *order.Order
	orderView   *order.OrderView
}

func OrderUpdateHandler(services *app.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &orderUpdateHandler{
			connections: services,
			params:      &order.OrderUpdateParams{},
			order:       &order.Order{},
			orderView:   &order.OrderView{},
		}
		h.Exec(c)
	}
}

func (h *orderUpdateHandler) Exec(c *gin.Context) {
	// prepare and check params
	// load the order from the database
	if err := h.prepare(c); err != nil {
		app.AbortWithErrorResponse(c, app.ErrFailedToLoadParams, err)
		return
	}

	if h.params.Status != nil {
		// status update handling
		if *h.params.Status == order.StatusCancelled {
			if h.order.Status == order.StatusCancelled.String() {
				c.JSON(http.StatusNoContent, nil) // no change
				return
			}
			if err := h.cancelOrder(c); err != nil { // allowed update
				app.AbortWithErrorResponse(c, app.ErrServerError, err)
				return
			}
		} else { // a valid enum is passed, but only cancelled can be set manually from this endpoint
			app.AbortWithErrorResponse(c, order.ErrInvalidOrderStatusUpdate, order.ErrInvalidOrderStatusUpdate.Err)
		}
	} else if h.params.ShippingAddress != nil {
		// shipping address update handling
		switch *h.params.ShippingAddress {
		case "":
			app.AbortWithErrorResponse(c, order.ErrShippingAddressRequired, order.ErrShippingAddressRequired.Err)
			return
		case h.order.ShippingAddress:
			c.JSON(http.StatusNoContent, nil) // no change
			return
		default:
			// update the order with the new shipping address
			h.order.ShippingAddress = *h.params.ShippingAddress
			h.order.UpdatedAt = time.Now()
			query := ordersql.NewUpdateOrderQuery(h.connections.DB, h.connections.Redis)
			query.Values.Order = h.order
			if err := query.Run(c); err != nil {
				app.AbortWithErrorResponse(c, app.ErrServerError, err)
				return
			}
		}
	}
	// fetch order view to send as response
	if err := h.buildResponse(c); err != nil {
		app.AbortWithErrorResponse(c, app.ErrServerError, err)
		return
	}
	c.JSON(http.StatusOK, h.orderView)
}

//////////////////////////////////////////////////////////////////////////////////////
// helper functions

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

	h.params.Sanitize()

	// load order if it exists
	query := ordersql.NewSelectOrderByIDQuery(h.connections.DB, h.connections.Redis)
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

	h.order = query.Order
	return nil
}

func (h *orderUpdateHandler) cancelOrder(c *gin.Context) error {
	// fetch related order items
	query := ordersql.NewSelectOrderItemListByOrderIDQuery(h.connections.DB, h.connections.Redis)
	query.Where.OrderID = h.orderID
	if err := query.Run(c); err != nil {
		return app.ErrServerError.Err
	}

	var prodIDs []int
	prodQuantityMap := make(map[int]int)
	for _, item := range query.Items {
		prodIDs = append(prodIDs, item.ProductID)
		prodQuantityMap[item.ProductID] = item.Quantity
	}

	// start a transaction to update the order status and return the stock to the products
	tx, err := h.connections.DB.BeginTx(c, pgx.TxOptions{})
	if err != nil {
		return app.ErrServerError.Err
	}

	h.order.Status = order.StatusCancelled.String()
	h.order.UpdatedAt = time.Now()

	// update the order status
	query2 := ordersql.NewUpdateOrderQuery(tx, h.connections.Redis)
	query2.Values.Order = h.order
	if err := query2.Run(c); err != nil {
		tx.Rollback(c)
		return app.ErrServerError.Err
	}

	// fetch the product list and lock for update to avoid concurrency issues
	query3 := prodsql.NewSelectProductListByIDsForUpdateQuery(tx)
	query3.Where.IDs = prodIDs
	if err := query3.Run(c); err != nil {
		tx.Rollback(c)
		return app.ErrServerError.Err
	}

	// return stock to products
	productStock := prodsql.ProductStock{}
	productStockList := make([]prodsql.ProductStock, len(query3.Products))

	for _, item := range query3.Products {
		productStock.ID = item.ID
		productStock.Stock = item.Stock + prodQuantityMap[item.ID]
		productStockList = append(productStockList, productStock)
	}

	query4 := prodsql.NewUpdateProductsStockByIDsQuery(tx)
	query4.Values.ProductStockList = productStockList
	if err := query4.Run(c); err != nil {
		tx.Rollback(c)
		return app.ErrServerError.Err
	}

	if err := tx.Commit(c); err != nil {
		return app.ErrServerError.Err
	}
	return nil
}

func (h *orderUpdateHandler) buildResponse(c *gin.Context) error {
	// fetch the order view after the updates
	query := ordersql.NewSelectOrderViewByIDQuery(h.connections.DB, h.connections.Redis)
	query.Where.ID = h.orderID
	if err := query.Run(c); err != nil {
		return app.ErrServerError.Err
	}

	// fetch the order items view to build the order
	query2 := ordersql.NewSelectOrderItemViewListByOrderIDQuery(h.connections.DB)
	query2.Where.OrderID = h.orderID
	if err := query2.Run(c); err != nil {
		return app.ErrServerError.Err
	}
	query.OrderView.OrderItems = query2.Items
	for _, item := range query2.Items {
		query.OrderView.TotalPrice += item.ItemsPrice
	}

	// index the order in elasticsearch
	indexer := orderesrc.NewOrderIndexer(h.connections.ES, "orders")
	if err := indexer.Run(c, query.OrderView); err != nil {
		return app.ErrServerError.Err
	}

	h.orderView = query.OrderView
	return nil
}
