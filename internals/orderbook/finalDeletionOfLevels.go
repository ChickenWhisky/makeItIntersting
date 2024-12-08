package orderbook

func (ob *OrderBook) FinalLevelDeletion() {
	for !ob.AsksLevelByLevel.Empty() {
		l, _ := ob.AsksLevelByLevel.Peek()
		Level := l.(*LevelBook)
		_, isLevelToBeDeleted := ob.ToBeDeletedLevels[Level.LevelID]
		if isLevelToBeDeleted {
			delete(ob.ToBeDeletedLevels, Level.LevelID)
			ob.AsksLevelByLevel.Dequeue()
		}
	}
	for !ob.BidsLevelByLevel.Empty() {
		l, _ := ob.BidsLevelByLevel.Peek()
		Level := l.(*LevelBook)
		_, isLevelToBeDeleted := ob.ToBeDeletedLevels[Level.LevelID]
		if isLevelToBeDeleted {
			delete(ob.ToBeDeletedLevels, Level.LevelID)
			ob.BidsLevelByLevel.Dequeue()
		}
	}
	for !ob.LimitAsksLevelByLevel.Empty() {
		l, _ := ob.LimitAsksLevelByLevel.Peek()
		Level := l.(*LevelBook)
		_, isLevelToBeDeleted := ob.ToBeDeletedLevels[Level.LevelID]
		if isLevelToBeDeleted {
			delete(ob.ToBeDeletedLevels, Level.LevelID)
			ob.LimitAsksLevelByLevel.Dequeue()
		}
	}
	for !ob.LimitBidsLevelByLevel.Empty() {
		l, _ := ob.LimitBidsLevelByLevel.Peek()
		Level := l.(*LevelBook)
		_, isLevelToBeDeleted := ob.ToBeDeletedLevels[Level.LevelID]
		if isLevelToBeDeleted {
			delete(ob.ToBeDeletedLevels, Level.LevelID)
			ob.LimitBidsLevelByLevel.Dequeue()
		}
	}

}
