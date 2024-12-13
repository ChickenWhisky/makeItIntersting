package orderbook

// DeleteLevelBook deletes a given Level Book
func (ob *OrderBook) DeleteLevelBook(lb *LevelBook) {
	if lb == nil {
		return // Return early if lb is nil
	}

	if lb.Orders != nil {
		lb.Orders.Clear() // Ensure Clear() is a valid method
	}

	if lb.Contracts != nil {
		for key := range lb.Contracts {
			delete(lb.Contracts, key)
		}
	}

	lb.Orders = nil
	lb.Contracts = nil
	ob.ToBeDeletedLevels[lb.LevelID] = lb
}
