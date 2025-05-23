package api

import (
	"github.com/timur-raja/order-tracking-rest-go/app"
	"github.com/timur-raja/order-tracking-rest-go/app/order/orderapi"
	"github.com/timur-raja/order-tracking-rest-go/app/user/userapi"

	"github.com/gin-gonic/gin"
)

// endpoints.go setsup the API endpoints for the application.

type Router struct {
	*app.Services
}

func NewRouter(services *app.Services) *Router {
	return &Router{
		Services: services,
	}
}

func (r *Router) Setup(router *gin.Engine) {
	public := router.Group("/")
	{
		public.POST("/signup", userapi.UserCreateHandler(r.Services))
		public.POST("/signin", userapi.UserSigninHandler(r.Services))
	}

	protected := router.Group("/")
	protected.Use(SessionAuth(r.Services.DB))
	{
		protected.POST(
			"/orders", orderapi.OrderCreateHandler(r.Services))
		protected.PATCH(
			"/orders/:order_id", orderapi.OrderUpdateHandler(r.Services)) // handles edit and cancellation
		protected.GET(
			"/orders", orderapi.OrderListHandler(r.Services))
		protected.GET(
			"/orders/:order_id", orderapi.OrderReadHandler(r.Services))
	}
}
