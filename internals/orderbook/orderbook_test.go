package orderbook

import (
	"github.com/ChickenWhisky/makeItIntersting/pkg/models"
	"testing"
)

func TestAddContract(t *testing.T) {
	ob := NewOrderBook()
	contract := models.Contract{
		OrderType:   "buy",
		RequestType: "add",
		Quantity:    10,
		Price:       100,
	}

	err := ob.AddContract(contract)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(ob.Orders) != 1 {
		t.Errorf("expected 1 contract, got %d", len(ob.Orders))
	}

}

func TestDeleteContract(t *testing.T) {
	ob := NewOrderBook()
	contract := models.Contract{
		ContractID:  "123",
		OrderType:   "buy",
		RequestType: "delete",
		Quantity:    10,
		Price:       100,
		UserID:      "user1",
	}

	ob.AddContract(contract)
	contract.RequestType = "delete"
	err := ob.CancelContract(contract)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(ob.Orders) != 0 {
		t.Errorf("expected 0 contracts, got %d", len(ob.Orders))
	}
}

func TestModifyContract(t *testing.T) {
	ob := NewOrderBook()
	contract := models.Contract{
		ContractID:  "123",
		OrderType:   "buy",
		RequestType: "add",
		Quantity:    10,
		Price:       100,
		UserID:      "user1",
	}

	ob.AddContract(contract)
	modifiedContract := models.Contract{
		ContractID: "123",
		OrderType:  "buy",
		Quantity:   20,
		Price:      150,
		UserID:     "user1",
	}

	ob.ModifyContract(modifiedContract)
	if ob.Orders["123"].Quantity != 20 {
		t.Errorf("expected quantity 20, got %d", ob.Orders["123"].Quantity)
	}

	if ob.Orders["123"].Price != 150 {
		t.Errorf("expected price 150, got %f", ob.Orders["123"].Price)
	}
}

func TestMatchOrders(t *testing.T) {
	ob := NewOrderBook()
	askContract := models.Contract{
		ContractID:  "ask1",
		OrderType:   "sell",
		RequestType: "add",
		Quantity:    10,
		Price:       100,
	}
	bidContract := models.Contract{
		ContractID:  "bid1",
		OrderType:   "buy",
		RequestType: "add",
		Quantity:    10,
		Price:       100,
	}

	ob.AddContract(askContract)
	ob.AddContract(bidContract)
	ob.MatchOrders()

	if len(ob.LastMatchedPrices) != 1 {
		t.Errorf("expected 1 matched trade, got %d", len(ob.LastMatchedPrices))
	}

	if ob.LastMatchedPrices[0].Price != 100 {
		t.Errorf("expected matched price 100, got %f", ob.LastMatchedPrices[0].Price)
	}
}
