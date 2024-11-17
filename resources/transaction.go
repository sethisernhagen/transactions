package resources

import (
	"errors"
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
		return errors.New("missing required Transaction fields")
	}
	if t.Transaction.Amount == 0 {
		return errors.New("missing required Amount field")
	}
	if !models.IsValidOperationTypeID(t.OperationType) {
		return errors.New("invalid OperationTypeID value")
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
