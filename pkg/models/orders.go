package models

// Order struct represents an order (buy/sell)
type Order struct {
	OrderID     string  `json:"order_id"`
	EventID     string  `json:"event_id"`
	SubEventID  string  `json:"subevent_id"`
	UserID      string  `json:"user_id"`
	RequestType string  `json:"request_type"`
	OrderType   string  `json:"order_type"`
	Price       float32 `json:"price" binding:"required,gt=0"`
	Quantity    int64   `json:"quantity" binding:"required,gte=1"`
	Timestamp   int64   `json:"timestamp"`
}

// NewOrder creates a new instance of Order
func NewOrder() *Order {
	return &Order{}
}

// SetOrderID sets the Order ID of the Order
func (o *Order) SetOrderID(ID string) {
	o.OrderID = ID
}

// GetOrderID returns the Order ID of the Order
func (o *Order) GetOrderID() string {
	return o.OrderID
}

// SetUserID sets the user ID of the Order
func (o *Order) SetUserID(ID string) {
	o.UserID = ID
}

// GetUserID returns the user ID of the Order
func (o *Order) GetUserID() string {
	return o.UserID
}

// SetRequestType sets the request type of the Order
func (o *Order) SetRequestType(rt string) {
	o.RequestType = rt
}

// GetRequestType returns the request type of the Order
func (o *Order) GetRequestType() string {
	return o.RequestType
}

// SetOrderType sets the order type of the Order
func (o *Order) SetOrderType(orderType string) {
	o.OrderType = orderType
}

// GetOrderType returns the order type of the Order
func (o *Order) GetOrderType() string {
	return o.OrderType
}

// SetPrice sets the price of the Order
func (o *Order) SetPrice(price float32) {
	o.Price = price
}

// GetPrice returns the price of the Order
func (o *Order) GetPrice() float32 {
	return o.Price
}

// SetQuantity sets the quantity of the Order
func (o *Order) SetQuantity(quantity int64) {
	o.Quantity = quantity
}

// GetQuantity returns the quantity of the Order
func (o *Order) GetQuantity() int64 {
	return o.Quantity
}

// SetTimestamp sets the timestamp of the Order
func (o *Order) SetTimestamp(timestamp int64) {
	o.Timestamp = timestamp
}

// GetTimestamp returns the timestamp of the Order
func (o *Order) GetTimestamp() int64 {
	return o.Timestamp
}

// SetEventID sets the event ID of the Order
func (o *Order) SetEventID(event string) {
	o.EventID = event
}

// GetEventID returns the event ID of the Order
func (o *Order) GetEventID() string {
	return o.EventID
}

// SetSubEventID sets the subevent ID of the Order
func (o *Order) SetSubEventID(event string) {
	o.SubEventID = event
}

// GetSubEventID returns the subevent ID of the Order
func (o *Order) GetSubEventID() string {
	return o.SubEventID
}
