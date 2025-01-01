package orderbook

import (
	"encoding/json"
	"github.com/ChickenWhisky/makeItIntersting/pkg/models"
	"github.com/charmbracelet/log"
	"os"
	"strconv"
	"time"
)

func (ob *OrderBook) LogHandler(lowestAskOrder *models.Order, highestBidOrder *models.Order) {
	// Log the trade
	trade := models.Trade{
		TradeID:          strconv.Itoa(ob.TradeNo),
		EventID:		  lowestAskOrder.EventID,
		SellerUserID:     lowestAskOrder.UserID,
		SellerOrderID: lowestAskOrder.OrderID,
		BuyerUserID:      highestBidOrder.UserID,
		BuyerOrderID:  highestBidOrder.OrderID,
		Price:            lowestAskOrder.Price,
		Quantity:         min(lowestAskOrder.Quantity, highestBidOrder.Quantity),
		Timestamp:        time.Now().UnixMilli(),
	}
	ob.TradeNo++
	ob.LastMatchedPrices = append(ob.LastMatchedPrices, trade)
	file, err := os.OpenFile("trades.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening file:", err)
		return
	}
	defer file.Close()

	tradeJSON, err := json.Marshal(trade)
	if err != nil {
		log.Printf("Error marshalling trade to JSON:", err)
		return
	}

	if _, err := file.Write(tradeJSON); err != nil {
		log.Printf("Error writing trade to file:", err)
		return
	}

	if _, err := file.WriteString("\n"); err != nil {
		log.Printf("Error writing newline to file:", err)
		return
	}
}
