package orderapi

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type orderListHandler struct {
	db *pgxpool.Pool
}

func OrderListHandler(db *pgxpool.Pool) gin.HandlerFunc {
	h := &orderListHandler{db: db}
	return h.Exec
}

func (h *orderListHandler) Exec(c *gin.Context) {
	return
}
