package orderapi

import (
	"net/http"

	"github.com/timur-raja/order-tracking-rest-go/app"
	"github.com/timur-raja/order-tracking-rest-go/app/order/orderesrc"

	"github.com/gin-gonic/gin"
)

type orderListHandler struct {
	connections *app.Services
}

func OrderListHandler(services *app.Services) gin.HandlerFunc {
	h := &orderListHandler{
		connections: services,
	}
	return h.Exec
}

func (h *orderListHandler) Exec(c *gin.Context) {
	esQuery := orderesrc.NewOrdersSearchQuery(h.connections.ES, "orders")
	if err := c.ShouldBindQuery(&esQuery.Params); err != nil {
		app.AbortWithErrorResponse(c, app.ErrFailedToLoadParams, err)
		return
	}

	if err := esQuery.Run(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, esQuery.Result.OrderViewList)
}
