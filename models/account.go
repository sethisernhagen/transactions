package models

import "time"

type Account struct {
	AccountID      int64     `json:"id,omitempty"`
	DocumentNumber string    `json:"document_number,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
}
