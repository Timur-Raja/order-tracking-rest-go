package userapi

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type userReadHandlerStruct struct {
	db *pgxpool.Pool
}

func UserReadHandler(db *pgxpool.Pool) gin.HandlerFunc {
	h := &userReadHandlerStruct{db: db}
	return h.Exec
}
func (h *userReadHandlerStruct) Exec(c *gin.Context) {
	return
}
