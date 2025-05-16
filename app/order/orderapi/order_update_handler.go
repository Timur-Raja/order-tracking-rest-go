package orderapi

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type orderUpdateHandler struct {
	db *pgxpool.Pool
}

func OrderUpdateHandler(db *pgxpool.Pool) gin.HandlerFunc {
	h := &orderUpdateHandler{db: db}
	return h.Exec
}

func (h *orderUpdateHandler) Exec(c *gin.Context) {
	return
}
