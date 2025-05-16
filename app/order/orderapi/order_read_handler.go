package orderapi

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type orderReadHandler struct {
	db *pgxpool.Pool
}

func OrderReadHandler(db *pgxpool.Pool) gin.HandlerFunc {
	h := &orderReadHandler{db: db}
	return h.Exec
}

func (h *orderReadHandler) Exec(c *gin.Context) {
	return
}
