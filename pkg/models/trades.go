package models

// Trade struct represents a completed Transaction and its details
type Trade struct {
	TradeID       string  `json:"trade_id"`
	EventID       string  `json:"event_id"`
	SubEventID    string  `json:"sub_event_id"`
	SellerUserID  string  `json:"seller_user_id"`
	SellerOrderID string  `json:"seller_Order_id"`
	BuyerUserID   string  `json:"buyer_user_id"`
	BuyerOrderID  string  `json:"buyer_Order_id"`
	Price         float32 `json:"price" binding:"required,gt=0"`
	Quantity      int64   `json:"quantity" binding:"required,gte=1"`
	Timestamp     int64   `json:"timestamp"`
}
