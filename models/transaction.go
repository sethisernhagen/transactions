package models

import "time"

type Transaction struct {
	TransactionID   int64     `json:"id"`
	AccountID       int64     `json:"account_id"`
	OperationTypeID int       `json:"operation_type_id"`
	Amount          float64   `json:"amount"`
	CreatedAt       time.Time `json:"created_at"`
}
