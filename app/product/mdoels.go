package product

type Product struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Price     float32 `json:"price"`
	Stock     int     `json:"stock"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt string  `json:"deleted_at"`
}
