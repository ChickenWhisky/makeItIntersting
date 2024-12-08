package orderbook

// DeleteLevelBook
func (ob *OrderBook) DeleteLevelBook(lb *LevelBook) {
	lb.Orders.Clear()

	if lb.Contracts != nil {
		for key := range lb.Contracts {
			delete(lb.Contracts, key)
		}
	}

	lb.Orders = nil
	lb.Contracts = nil

}
