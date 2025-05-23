package orderapi

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/timur-raja/order-tracking-rest-go/app"
	"github.com/timur-raja/order-tracking-rest-go/app/order"
	"github.com/timur-raja/order-tracking-rest-go/app/order/ordersql"
)

type orderReadHandler struct {
	connections *app.Services
}

func OrderReadHandler(services *app.Services) gin.HandlerFunc {
	h := &orderReadHandler{
		connections: services,
	}
	return h.exec
}

func (h *orderReadHandler) exec(c *gin.Context) {
	// fetch id from the URL
	idParam := c.Param("order_id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		app.AbortWithErrorResponse(c, app.ErrFailedToLoadParams, err)
		return
	}

	query := ordersql.NewSelectOrderViewByIDQuery(h.connections.DB, h.connections.Redis)
	query.Where.ID = id
	if err := query.Run(c); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			app.AbortWithErrorResponse(c, order.ErrOrderNotFound, err)
		} else {
			app.AbortWithErrorResponse(c, app.ErrServerError, err)
		}
		return
	}

	query2 := ordersql.NewSelectOrderItemViewListByOrderIDQuery(h.connections.DB)
	query2.Where.OrderID = id
	if err := query2.Run(c); err != nil {
		app.AbortWithErrorResponse(c, app.ErrServerError, err)
		return
	}
	query.OrderView.OrderItems = query2.Items

	c.JSON(200, query.OrderView)
}
