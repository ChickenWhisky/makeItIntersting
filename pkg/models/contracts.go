package models

// Contract struct represents an order (buy/sell)
type Contract struct {
	ContractID  string  `json:"contract_id"`
	UserID      string  `json:"user_id"`
	RequestType string  `json:"request_type"`
	OrderType   string  `json:"order_type"`
	Price       float32 `json:"price" binding:"required,gt=0"`
	Quantity    int64   `json:"quantity" binding:"required,gte=1"`
	Timestamp   int64   `json:"timestamp"`
}

// NewContract creates a new instance of Contract
func NewContract() *Contract {
	return &Contract{}
}

// SetContractID sets the contract ID of the contract
func (c *Contract) SetContractID(ID string) {
	c.ContractID = ID
}

// GetContractID returns the contract ID of the contract
func (c *Contract) GetContractID() string {
	return c.ContractID
}

// SetUserID sets the user ID of the contract
func (c *Contract) SetUserID(ID string) {
	c.UserID = ID
}

// GetUserID returns the user ID of the contract
func (c *Contract) GetUserID() string {
	return c.UserID
}

// SetRequestType sets the request type of the contract
func (c *Contract) SetRequestType(rt string) {
	c.RequestType = rt
}

// GetRequestType returns the request type of the contract
func (c *Contract) GetRequestType() string {
	return c.RequestType
}

// SetOrderType sets the order type of the contract
func (c *Contract) SetOrderType(orderType string) {
	c.OrderType = orderType
}

// GetOrderType returns the order type of the contract
func (c *Contract) GetOrderType() string {
	return c.OrderType
}

// SetPrice sets the price of the contract
func (c *Contract) SetPrice(price float32) {
	c.Price = price
}

// GetPrice returns the price of the contract
func (c *Contract) GetPrice() float32 {
	return c.Price
}

// SetQuantity sets the quantity of the contract
func (c *Contract) SetQuantity(quantity int64) {
	c.Quantity = quantity
}

// GetQuantity returns the quantity of the contract
func (c *Contract) GetQuantity() int64 {
	return c.Quantity
}

// SetTimestamp sets the timestamp of the contract
func (c *Contract) SetTimestamp(timestamp int64) {
	c.Timestamp = timestamp
}

// GetTimestamp returns the timestamp of the contract
func (c *Contract) GetTimestamp() int64 {
	return c.Timestamp
}
