package models

type Contract struct {
	UserID    string  `json:"user_id"`
	OrderType string  `json:"order_type"`
	Price     float64 `json:"price" binding:"required,gt=0"`
	Quantity  int     `json:"quantity" binding:"required,gte=1"`
	Timestamp int64   `json:"timestamp"`
}
