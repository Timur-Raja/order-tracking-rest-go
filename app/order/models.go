package order

import (
	"time"
)

type Order struct {
	ID              int        `db:"id"`
	UserID          int        `db:"user_id"`
	Status          string     `db:"status"`
	ShippingAddress string     `db:"shipping_address"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at"`
}

type OrderItem struct {
	ID        int     `db:"id"`
	OrderID   int     `db:"order_id"`
	ProductID int     `db:"product_id"`
	Quantity  int     `db:"quantity"`
	Price     float32 `db:"price"`
}
