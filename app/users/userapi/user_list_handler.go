package userapi

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type userListHandler struct {
	db *pgxpool.Pool
}

func UserListHandler(db *pgxpool.Pool) gin.HandlerFunc {
	h := &userListHandler{db: db}
	return h.Exec
}

func (h *userListHandler) Exec(c *gin.Context) {
	c.JSON(200, gin.H{"message": "called users list"})
}
