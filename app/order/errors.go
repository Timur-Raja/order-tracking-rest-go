package order

import (
	"net/http"

	"github.com/timur-raja/order-tracking-rest-go/app"
)

var (
	ErrOrderNotFound           = app.NewAPIError(http.StatusNotFound, "order not found")
	ErrShippingAddressRequired = app.NewAPIError(http.StatusBadRequest, "shipping address is required")
)
