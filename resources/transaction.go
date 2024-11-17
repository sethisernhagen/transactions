package resources

import (
	"errors"
	"fmt"
	"net/http"
	"transactions/models"
)

type TransactionResponse struct {
	*models.Transaction
	OperationType string `json:"operation_type"`
}

type NewTransactionRequest struct {
	*models.Transaction
	OperationType            string `json:"operation_type"`
	ProtectedOperationTypeID string `json:"operation_type_id"` // override
	ProtectedID              string `json:"id"`                // override
}

func (t *NewTransactionRequest) Bind(r *http.Request) error {
	if t.Transaction == nil {
		return errors.New("Missing required Transaction fields")
	}
	if t.Transaction.Amount == 0 {
		return errors.New("Missing required Amount field")
	}
	if !models.IsValidOperationTypeID(t.OperationType) {
		return errors.New("Invalid OperationTypeID value")
	}
	fmt.Println(t.OperationType)
	for _, validType := range models.NegativeOperationTypeIDs {
		if t.OperationType == validType {
			if t.Transaction.Amount > 0 {
				return errors.New("Amount must be negative for Purchases or Withdrawals")
			}
		}
	}

	for _, validType := range models.PositiveOperationTypeIDs {
		if t.OperationType == validType {
			if t.Transaction.Amount < 0 {
				return errors.New("Amount must be positive for CreditVouchers")
			}
		}
	}

	t.Transaction.OperationTypeID = t.Transaction.OperationTypeID.GetID(t.OperationType)
	t.ProtectedID = ""
	return nil
}

func NewTransactionResponse(transaction *models.Transaction) *TransactionResponse {
	resp := &TransactionResponse{Transaction: transaction}
	return resp
}

func (tr *TransactionResponse) Render(w http.ResponseWriter, r *http.Request) error {
	tr.OperationType = tr.Transaction.OperationTypeID.String()
	return nil
}
