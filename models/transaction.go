package models

import "time"

type Transaction struct {
	TransactionID   int64
	AccountID       int64
	OperationTypeID int
	Amount          float64
	CreatedAt       time.Time
}
