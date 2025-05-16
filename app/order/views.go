package order

type OrderView struct {
	ID           int             `json:"id"`
	UserID       int             `json:"userID"`
	UserFullName string          `json:"userFUllName"`
	UserEmail    string          `json:"userEmail"`
	CreatedAt    string          `json:"createdAt"`
	UpdatedAt    string          `json:"updatedAt"`
	DeletedAt    string          `json:"deletedAt"`
	TotalPrice   float64         `json:"totalPrice"`
	OrderItems   []OrderItemView `json:"orderItems"`
}

type OrderItemView struct {
	ID           int     `json:"id"`
	OrderID      int     `json:"orderID"`
	ProductID    int     `json:"productID"`
	ProductName  string  `json:"productName"`
	Quantity     int     `json:"quantity"`
	ProductPrice float64 `json:"productPrice"`
}
