package orderbook

// DeleteLevelBook deletes a given Level Book
func (ob *OrderBook) DeleteLevelBook(lb *LevelBook) {
	lb.Orders.Clear()

	if lb.Contracts != nil {
		for key := range lb.Contracts {
			delete(lb.Contracts, key)
		}
	}
	lb.Orders = nil
	lb.Contracts = nil
	ob.ToBeDeletedLevels[lb.LevelID] = lb

}
