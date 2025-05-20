package order

type OrderCreateParams struct {
	OrderItems      []OrderItemCreateParams `json:"orderItems" binding:"required"`
	ShippingAddress string                  `json:"shippingAddress" binding:"required"`
}

type OrderItemCreateParams struct {
	ProductID int `json:"productID" binding:"required"`
	Quantity  int `json:"quantity" binding:"required"`
}

// update params are pointers to allow partial updates
// and to differentiate between zero values and unset values
type OrderUpdateParams struct {
	Status          *string                 `json:"status"`
	OrderItems      []OrderItemUpdateParams `json:"orderItems"`
	ShippingAddress *string                 `json:"shippingAddress"`
}

type OrderItemUpdateParams struct {
	Quantity  *int `json:"quantity" `
	ProductID *int `json:"productID"`
}
