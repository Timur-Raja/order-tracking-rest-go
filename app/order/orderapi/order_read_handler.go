package orderapi

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/timur-raja/order-tracking-rest-go/app"
	"github.com/timur-raja/order-tracking-rest-go/app/order/ordersql"
)

type orderReadHandler struct {
	db *pgxpool.Pool
}

func OrderReadHandler(db *pgxpool.Pool) gin.HandlerFunc {
	return (&orderReadHandler{db: db}).exec
}

func (h *orderReadHandler) exec(c *gin.Context) {
	// fetch id from the URL
	idParam := c.Param("order_id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		app.AbortWithErrorResponse(c, app.ErrFailedToLoadParams, err)
		return
	}

	query := ordersql.NewSelectOrderViewByIDQuery(h.db)
	query.Where.ID = id
	if err := query.Run(c); err != nil {
		if err == pgx.ErrNoRows {
			app.AbortWithErrorResponse(c, app.ErrResourceNotFound, err)
		} else {
			app.AbortWithErrorResponse(c, app.ErrServerError, err)
		}
		return
	}

	c.JSON(200, query.OrderView)
}
