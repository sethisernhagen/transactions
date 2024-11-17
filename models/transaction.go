package models

import "time"

type Transaction struct {
	TransactionID   int64           `json:"id,omitempty"`
	AccountID       int64           `json:"account_id,omitempty"`
	OperationTypeID OperationTypeID `json:"operation_type_id,omitempty"`
	Amount          float64         `json:"amount,omitempty"`
	CreatedAt       time.Time       `json:"created_at,omitempty"`
}

type OperationTypeID int

const (
	Purchase             OperationTypeID = 1
	PurchaseInstallments OperationTypeID = 2
	Withdrawal           OperationTypeID = 3
	CreditVoucher        OperationTypeID = 4
)

var ValidOperationTypeIDs = []string{"Purchase", "PurchaseInstallments", "Withdrawal", "CreditVoucher"}
var NegativeOperationTypeIDs = []string{"Purchase", "PurchaseInstallments", "Withdrawal"}
var PositiveOperationTypeIDs = []string{"CreditVoucher"}

func (o OperationTypeID) String() string {
	switch o {
	case Purchase:
		return "Purchase"
	case PurchaseInstallments:
		return "PurchaseInstallments"
	case Withdrawal:
		return "Withdrawal"
	case CreditVoucher:
		return "CreditVoucher"
	default:
		return "Unknown"
	}
}
func (o OperationTypeID) GetID(name string) OperationTypeID {
	switch name {
	case "Purchase":
		return Purchase
	case "PurchaseInstallments":
		return PurchaseInstallments
	case "Withdrawal":
		return Withdrawal
	case "CreditVoucher":
		return CreditVoucher
	default:
		return -1
	}
}

func IsValidOperationTypeID(name string) bool {
	switch name {
	case "Purchase":
		return true
	case "PurchaseInstallments":
		return true
	case "Withdrawal":
		return true
	case "CreditVoucher":
		return true
	default:
		return false
	}
}
