package order

type OrderCreateParams struct {
	OrderItems []OrderItemCreateParams `json:"orderItems"`
}

type OrderItemCreateParams struct {
	ProductID int `json:"productID"`
	Quantity  int `json:"quantity"`
}
