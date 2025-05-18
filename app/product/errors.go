package product

import (
	"net/http"

	"github.com/timur-raja/order-tracking-rest-go/app"
)

var (
	ErrProductsNotFound = app.NewAPIError(http.StatusNotFound, "one or more products not found")
)
