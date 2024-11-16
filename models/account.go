package models

import "time"

type Account struct {
	AccountID      int64     `json:"id"`
	DocumentNumber string    `json:"document_number"`
	CreatedAt      time.Time `json:"created_at"`
}
