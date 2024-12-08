package models

// Trade struct represents a completed Transaction and its details
type Trade struct {
	TradeID          string  `json:"contract_id"`
	SellerUserID     string  `json:"seller_user_id"`
	SellerContractID string  `json:"seller_contract_id"`
	BuyerUserID      string  `json:"buyer_user_id"`
	BuyerContractID  string  `json:"buyer_contract_id"`
	Price            float32 `json:"price" binding:"required,gt=0"`
	Quantity         int64   `json:"quantity" binding:"required,gte=1"`
	Timestamp        int64   `json:"timestamp"`
}
