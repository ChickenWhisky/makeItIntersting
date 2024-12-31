package orderbook

import (
	"encoding/json"
	"github.com/ChickenWhisky/makeItIntersting/pkg/models"
	"github.com/charmbracelet/log"
	"os"
	"strconv"
	"time"
)

func (ob *OrderBook) LogHandler(lowestAskContract *models.Contract, highestBidContract *models.Contract) {
	// Log the trade
	trade := models.Trade{
		TradeID:          strconv.Itoa(ob.TradeNo),
		SellerUserID:     lowestAskContract.UserID,
		SellerContractID: lowestAskContract.ContractID,
		BuyerUserID:      highestBidContract.UserID,
		BuyerContractID:  highestBidContract.ContractID,
		Price:            lowestAskContract.Price,
		Quantity:         min(lowestAskContract.Quantity, highestBidContract.Quantity),
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
