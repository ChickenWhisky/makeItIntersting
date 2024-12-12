package orderbook

import (
	"encoding/json"
	"github.com/ChickenWhisky/makeItIntersting/pkg/models"
	"log"
	"os"
	"time"
)

func logHandler(lowestAskContract *models.Contract, highestBidContract *models.Contract) {
	// Log the trade
	trade := models.Trade{
		TradeID:          "some_unique_trade_id", // You need to generate a unique ID
		SellerUserID:     lowestAskContract.UserID,
		SellerContractID: lowestAskContract.ContractID,
		BuyerUserID:      highestBidContract.UserID,
		BuyerContractID:  highestBidContract.ContractID,
		Price:            lowestAskContract.Price,
		Quantity:         min(lowestAskContract.Quantity, highestBidContract.Quantity),
		Timestamp:        time.Now().UnixMilli(),
	}

	file, err := os.OpenFile("trades.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	tradeJSON, err := json.Marshal(trade)
	if err != nil {
		log.Println("Error marshalling trade to JSON:", err)
		return
	}

	if _, err := file.Write(tradeJSON); err != nil {
		log.Println("Error writing trade to file:", err)
		return
	}

	if _, err := file.WriteString("\n"); err != nil {
		log.Println("Error writing newline to file:", err)
		return
	}
}
