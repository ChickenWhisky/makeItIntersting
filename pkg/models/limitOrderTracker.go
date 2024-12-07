package models

type LimitOrderTracker struct {
	Price     float32    `json:"Price of Contract to be traded"`
	Contracts []Contract `json:"contract_id"`
}
