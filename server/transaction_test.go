package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"transactions/models"
	"transactions/resources"
	"transactions/stores"

	"github.com/stretchr/testify/assert"
)

func TestTransaction_Create(t *testing.T) {
	s := NewServer()
	_, err := stores.GetTestDB(t)
	assert.NoError(t, err)

	ts := httptest.NewServer(s.Router)
	defer ts.Close()
	acct := createAccountFixture(t, ts)

	res, err := http.Post(
		ts.URL+"/transaction", "application/json",
		bytes.NewBufferString(fmt.Sprintf(`{"account_id":%d,"operation_type": "Purchase","amount":-100.0}`, acct.AccountID)),
	)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	resBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	transaction := resources.TransactionResponse{}
	assert.NoError(t, json.Unmarshal(resBody, &transaction))
	assert.IsType(t, int64(0), transaction.TransactionID)
	assert.Equal(t, acct.AccountID, transaction.AccountID)
	assert.Equal(t, models.Purchase.String(), transaction.OperationType)
	assert.Equal(t, -100.0, transaction.Amount)
	assert.True(t, transaction.CreatedAt.After(time.Unix(0, 0)), "transaction.CreatedAt should be after Unix epoch, got %v", transaction.CreatedAt)
}

func TestTransaction_Retrieve(t *testing.T) {
	s := NewServer()
	_, err := stores.GetTestDB(t)
	assert.NoError(t, err)
	ts := httptest.NewServer(s.Router)
	defer ts.Close()
	acct := createAccountFixture(t, ts)
	trx := createTransactionFixture(t, ts, acct.AccountID)

	res, err := http.Get(ts.URL + "/transaction/" + fmt.Sprint(trx.TransactionID))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	resBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	transaction := resources.TransactionResponse{}
	assert.NoError(t, json.Unmarshal(resBody, &transaction))
	assert.Equal(t, trx.AccountID, transaction.AccountID)
	assert.Equal(t, trx.OperationTypeID.String(), transaction.OperationType)
	assert.Equal(t, trx.Amount, transaction.Amount)
	assert.Equal(t, trx.CreatedAt, transaction.CreatedAt)
}

func TestTransaction_Retrieve_NotFound(t *testing.T) {
	s := NewServer()
	_, err := stores.GetTestDB(t)
	assert.NoError(t, err)
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/transaction/1")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestTransaction_Retrieve_BadRequest(t *testing.T) {
	s := NewServer()
	_, err := stores.GetTestDB(t)
	assert.NoError(t, err)
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/transaction/abc")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestTransaction_Create_BadRequest(t *testing.T) {
	s := NewServer()
	_, err := stores.GetTestDB(t)
	assert.NoError(t, err)
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	res, err := http.Post(ts.URL+"/transaction", "application/json", bytes.NewBufferString(`{}`))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestTransaction_Create_BadOperationType(t *testing.T) {
	s := NewServer()
	_, err := stores.GetTestDB(t)
	assert.NoError(t, err)
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	res, err := http.Post(ts.URL+"/transaction", "application/json", bytes.NewBufferString(`{"account_id":1,"operation_type":"Unknown","amount":100.0}`))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestTransaction_Create_BadPurchaseAmount(t *testing.T) {
	s := NewServer()
	_, err := stores.GetTestDB(t)
	assert.NoError(t, err)
	ts := httptest.NewServer(s.Router)
	defer ts.Close()
	acct := createAccountFixture(t, ts)

	res, err := http.Post(ts.URL+"/transaction", "application/json", bytes.NewBufferString(fmt.Sprintf(`{"account_id":%d,"operation_type":"Purchase","amount":100.0}`, acct.AccountID)))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Contains(t, getResponseBodyString(t, res), "Amount must be negative for Purchases or Withdrawals")
}
func TestTransaction_Create_BadOperationTypeAmount(t *testing.T) {
	s := NewServer()
	_, err := stores.GetTestDB(t)
	assert.NoError(t, err)
	ts := httptest.NewServer(s.Router)
	defer ts.Close()
	acct := createAccountFixture(t, ts)

	operationTypes := []string{"Purchase", "PurchaseInstallments", "Withdrawal"}
	for _, opType := range operationTypes {
		res, err := http.Post(
			ts.URL+"/transaction", "application/json",
			bytes.NewBufferString(fmt.Sprintf(`{"account_id":%d,"operation_type":"%s","amount":100.0}`, acct.AccountID, opType)),
		)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Contains(t, getResponseBodyString(t, res), "Amount must be negative for Purchases or Withdrawals")
	}
}

func TestTransaction_Create_BadOperationTypeAmountPositive(t *testing.T) {
	s := NewServer()
	_, err := stores.GetTestDB(t)
	assert.NoError(t, err)
	ts := httptest.NewServer(s.Router)
	defer ts.Close()
	acct := createAccountFixture(t, ts)

	operationTypes := []string{"CreditVoucher"}
	for _, opType := range operationTypes {
		res, err := http.Post(
			ts.URL+"/transaction", "application/json",
			bytes.NewBufferString(fmt.Sprintf(`{"account_id":%d,"operation_type":"%s","amount":-100.0}`, acct.AccountID, opType)),
		)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Contains(t, getResponseBodyString(t, res), "Amount must be positive for CreditVouchers")
	}
}

func TestTransaction_Create_PositiveOperationTypes(t *testing.T) {
	s := NewServer()
	_, err := stores.GetTestDB(t)
	assert.NoError(t, err)
	ts := httptest.NewServer(s.Router)
	defer ts.Close()
	acct := createAccountFixture(t, ts)

	operationTypes := []string{"CreditVoucher"}
	for _, opType := range operationTypes {
		res, err := http.Post(
			ts.URL+"/transaction", "application/json",
			bytes.NewBufferString(fmt.Sprintf(`{"account_id":%d,"operation_type":"%s","amount":100.0}`, acct.AccountID, opType)),
		)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode)
	}
}

func TestTransaction_Create_NegativeOperationTypes(t *testing.T) {
	s := NewServer()
	_, err := stores.GetTestDB(t)
	assert.NoError(t, err)
	ts := httptest.NewServer(s.Router)
	defer ts.Close()
	acct := createAccountFixture(t, ts)

	operationTypes := []string{"Purchase", "PurchaseInstallments", "Withdrawal"}
	for _, opType := range operationTypes {
		res, err := http.Post(
			ts.URL+"/transaction", "application/json",
			bytes.NewBufferString(fmt.Sprintf(`{"account_id":%d,"operation_type":"%s","amount":-100.0}`, acct.AccountID, opType)),
		)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode)
	}
}

func getResponseBodyString(t *testing.T, res *http.Response) string {
	resBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	return string(resBody)
}
