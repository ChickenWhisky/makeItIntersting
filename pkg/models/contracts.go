package models

// Contract struct represents an order (buy/sell)
type Contract struct {
	ContractID  string  `json:"contract_id"`
	UserID      string  `json:"user_id"`
	RequestType string  `json:"request_type"`
	OrderType   string  `json:"order_type"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Quantity    int     `json:"quantity" binding:"required,gte=1"`
	Timestamp   int64   `json:"timestamp"`
}
