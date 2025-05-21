package order

import "text/template"

type OrderCreateParams struct {
	OrderItems      []OrderItemCreateParams `json:"orderItems" binding:"required"`
	ShippingAddress string                  `json:"shippingAddress" binding:"required"`
}

// escape HTML to prevent XSS attacks
func (p *OrderCreateParams) Sanitize() {
	p.ShippingAddress = template.HTMLEscapeString(p.ShippingAddress)
}

type OrderItemCreateParams struct {
	ProductID int `json:"productID" binding:"required"`
	Quantity  int `json:"quantity" binding:"required"`
}

// update params are pointers to allow partial updates
// and to differentiate between zero values and unset values
type OrderUpdateParams struct {
	Status          *OrderStatus            `json:"status"` //custom enum type
	OrderItems      []OrderItemUpdateParams `json:"orderItems"`
	ShippingAddress *string                 `json:"shippingAddress"`
}

func (p *OrderUpdateParams) Sanitize() {
	if p.ShippingAddress != nil {
		*p.ShippingAddress = template.HTMLEscapeString(*p.ShippingAddress)
	}

}

type OrderItemUpdateParams struct {
	Quantity  *int `json:"quantity" `
	ProductID *int `json:"productID"`
}
