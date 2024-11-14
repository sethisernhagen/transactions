package models

import "time"

type Account struct {
	AccountID      int64
	DocumentNumber string
	CreatedAt      time.Time
}
