package resources

import (
	"errors"
	"net/http"
	"transactions/models"
)

type TransactionResponse struct {
	*models.Transaction
}

type NewTransactionRequest struct {
	*models.Transaction
	ProtectedID string `json:"id"` // override
}

func (t *NewTransactionRequest) Bind(r *http.Request) error {
	if t.Transaction == nil {
		return errors.New("missing required Transaction fields")
	}
	if t.Transaction.OperationTypeID == 0 {
		return errors.New("missing required OperationTypeID field")
	}
	if t.Transaction.Amount == 0 {
		return errors.New("missing required Amount field")
	}

	t.ProtectedID = ""
	return nil
}

func NewTransactionResponse(transaction *models.Transaction) *TransactionResponse {
	resp := &TransactionResponse{Transaction: transaction}
	return resp
}

func (rd *TransactionResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
