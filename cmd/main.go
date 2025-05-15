package main

import (
	"github.com/timur-raja/order-tracking-rest-go/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r)

	// Start the server
	if err := r.Run(); err != nil {
		panic(err)
	}
}
