package api

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
	public := router.Group("/")
	{
		public.POST("/signup", userapi.UserCreateHandler(r.DB))
		public.POST("/signin", userapi.UserSigninHandler(r.DB))
	}

	protected := router.Group("/")
	protected.Use(SessionAuth(r.DB))
	{
		protected.POST(
			"/orders", orderapi.OrderCreateHandler(r.DB))
		protected.PATCH(
			"/orders/:order_id", orderapi.OrderUpdateHandler(r.DB))
		protected.GET(
			"/orders", orderapi.OrderListHandler(r.DB))
		protected.GET(
			"/orders/:order_id", orderapi.OrderReadHandler(r.DB))
		protected.DELETE(
			"/orders/:order_id", orderapi.OrderDeleteHandler(r.DB))
	}
}
