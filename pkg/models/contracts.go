package models

// Contract struct represents an order (buy/sell)

type Contract struct {
	ContractID  string  `json:"contract_id"`
	UserID      string  `json:"user_id"`
	RequestType string  `json:"request_type (Request for deletion or addition of a contract)"`
	OrderType   string  `json:"order_type (buy/sell/limit_sell/limit_buy)"`
	Price       float32 `json:"price" binding:"required,gt=0"`
	Quantity    int64   `json:"quantity" binding:"required,gte=1"`
	Timestamp   int64   `json:"timestamp"`
}
