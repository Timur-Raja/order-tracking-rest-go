package app

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/timur-raja/order-tracking-rest-go/app/order/orderapi"
	"github.com/timur-raja/order-tracking-rest-go/app/user/userapi"

	"github.com/gin-gonic/gin"
)

// endpoints.go setsup the API endpoints for the application.

type Router struct {
	DB *pgxpool.Pool
}

func NewRouter(db *pgxpool.Pool) *Router {
	return &Router{DB: db}
}

func (r *Router) Setup(router *gin.Engine) {
	router.POST(
		"/signup", userapi.UserCreateHandler(r.DB))
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
