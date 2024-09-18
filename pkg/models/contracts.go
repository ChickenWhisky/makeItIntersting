package models

// Contract struct represents an order (buy/sell)
type Contract struct {
	UserID    string  `json:"user_id"`
	OrderType string  `json:"order_type"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Timestamp int64   `json:"timestamp"`
}
