package order

import (
	"time"
)

type Order struct {
	ID        int
	UserID    int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type OrderItem struct {
	OrderID   int
	ProductID int
	Quantity  int
	Price     float64
}
