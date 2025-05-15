package routes

import (
	"github.com/timur-raja/order-tracking-rest-go/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/users", handlers.HandleGet)
	router.POST("/orders", handlers.HandlePost)
}
