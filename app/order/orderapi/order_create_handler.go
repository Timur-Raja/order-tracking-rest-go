package orderapi

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type orderCreateHandler struct {
	db *pgxpool.Pool
}

func OrderCreateHandler(db *pgxpool.Pool) gin.HandlerFunc {
	h := &orderCreateHandler{db: db}
	return h.Exec
}

func (h *orderCreateHandler) Exec(c *gin.Context) {
	return
}
