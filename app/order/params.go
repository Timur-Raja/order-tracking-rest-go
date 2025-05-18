package order

type OrderCreateParams struct {
	OrderItems []OrderItemCreateParams
}

type OrderItemCreateParams struct {
	s
	ProductID int `json:"productID"`
	Quantity  int `json:"quantity"`
}
