package order

import "time"

type OrderView struct {
	ID         int              `json:"id" db:"id"`
	UserID     int              `json:"userID" db:"user_id"`
	Status     string           `json:"status" db:"status"`
	UserName   string           `json:"userName" db:"user_name"`
	UserEmail  string           `json:"userEmail" db:"user_email"`
	CreatedAt  time.Time        `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time        `json:"updatedAt" db:"updated_at"`
	DeletedAt  *time.Time       `json:"deletedAt" db:"deleted_at"`
	TotalPrice float32          `json:"totalPrice"`
	OrderItems []*OrderItemView `json:"orderItems"`
}

type OrderItemView struct {
	ID                    int     `json:"id" db:"id"`
	OrderID               int     `json:"orderID" db:"order_id"`
	ProductID             int     `json:"productID" db:"product_id"`
	ProductName           string  `json:"productName" db:"product_name"`
	Quantity              int     `json:"quantity" db:"ordered_quantity"`
	ItemsPrice            float32 `json:"itemsPrice" db:"items_price"`
	ProductPrice          float32 `json:"productPrice" db:"product_price"`
	RemainingProductStock int     `json:"productStock" db:"remaining_product_stock"`
}
