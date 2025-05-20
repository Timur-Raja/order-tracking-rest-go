package orderapi

import (
	"net/http"

	"github.com/timur-raja/order-tracking-rest-go/app"
	"github.com/timur-raja/order-tracking-rest-go/app/order/orderesrc"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/olivere/elastic/v7"
)

type orderListHandler struct {
	db *pgxpool.Pool
	es *elastic.Client
}

func OrderListHandler(db *pgxpool.Pool, es *elastic.Client) gin.HandlerFunc {
	h := &orderListHandler{
		db: db,
		es: es,
	}
	return h.Exec
}

func (h *orderListHandler) Exec(c *gin.Context) {
	esQuery := orderesrc.NewOrdersViewSearchQuery(h.es, "orders")
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
