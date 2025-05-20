package api

import (
	"github.com/timur-raja/order-tracking-rest-go/app/order/orderapi"
	"github.com/timur-raja/order-tracking-rest-go/app/user/userapi"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/olivere/elastic/v7"
)

// endpoints.go setsup the API endpoints for the application.

type Router struct {
	db *pgxpool.Pool
	es *elastic.Client
}

func NewRouter(db *pgxpool.Pool, es *elastic.Client) *Router {
	return &Router{
		db: db,
		es: es,
	}
}

func (r *Router) Setup(router *gin.Engine) {
	public := router.Group("/")
	{
		public.POST("/signup", userapi.UserCreateHandler(r.db))
		public.POST("/signin", userapi.UserSigninHandler(r.db))
	}

	protected := router.Group("/")
	protected.Use(SessionAuth(r.db))
	{
		protected.POST(
			"/orders", orderapi.OrderCreateHandler(r.db, r.es))
		protected.PATCH(
			"/orders/:order_id", orderapi.OrderUpdateHandler(r.db, r.es)) // handles edit and cancellation
		protected.GET(
			"/orders", orderapi.OrderListHandler(r.db, r.es))
		protected.GET(
			"/orders/:order_id", orderapi.OrderReadHandler(r.db))
	}
}
