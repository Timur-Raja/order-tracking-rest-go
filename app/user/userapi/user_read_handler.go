package userapi

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type userReadHandler struct {
	db *pgxpool.Pool
}

func UserReadHandler(db *pgxpool.Pool) gin.HandlerFunc {
	h := &userReadHandler{db: db}
	return h.Exec
}
func (h *userReadHandler) Exec(c *gin.Context) {
	return
}
