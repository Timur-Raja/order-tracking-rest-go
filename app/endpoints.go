package app

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/timur-raja/order-tracking-rest-go/app/order/orderapi"
	"github.com/timur-raja/order-tracking-rest-go/app/users/userapi"

	"github.com/gin-gonic/gin"
)

type Router struct {
	DB *pgxpool.Pool
}

func NewRouter(db *pgxpool.Pool) *Router {
	return &Router{DB: db}
}

func (r *Router) Setup(router *gin.Engine) {
	router.GET(
		"/users", userapi.UserListHandler(r.DB))
	router.GET(
		"/users/:user_id", userapi.UserReadHandler(r.DB))
	router.POST(
		"/orders", orderapi.OrderCreateHandler(r.DB))
	router.PATCH(
		"/orders/:order_id", orderapi.OrderUpdateHandler(r.DB))
	router.GET(
		"/orders", orderapi.OrderListHandler(r.DB))
	router.GET(
		"/orders/:order_id", orderapi.OrderReadHandler(r.DB))
	router.DELETE(
		"/orders/:order_id", orderapi.OrderDeleteHandler(r.DB))
}
