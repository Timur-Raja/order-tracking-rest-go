package orderapi

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type orderDeleteHandler struct {
	db *pgxpool.Pool
}

func OrderDeleteHandler(db *pgxpool.Pool) gin.HandlerFunc {
	h := &orderCreateHandler{db: db}
	return h.Exec
}

func (h *orderDeleteHandler) Exec(c *gin.Context) {
	return
}
