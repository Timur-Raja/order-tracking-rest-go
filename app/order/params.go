package order

type OrderCreateParams struct {
	OrderItems []OrderItemCreateParams
}

type OrderItemCreateParams struct {
	ProductID int `json:"productID"`
	Quantity  int `json:"quantity"`
}
