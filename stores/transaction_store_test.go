package stores

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionStore_Create(t *testing.T) {
	db, err := getTestDB(t)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	assert.NoError(t, err)
	accountStore := NewAccountStore(db)
	transactionStore := NewTransactionStore(db)

	acc, err := accountStore.Create(accountFixture1)
	assert.NoError(t, err)
	transactionFixture1.AccountID = acc.AccountID
	createdTransaction, err := transactionStore.Create(transactionFixture1)
	assert.NoError(t, err)
	assert.Equal(t, transactionFixture1, createdTransaction)
}

func TestTransactionStore_Retrieve(t *testing.T) {
	db, err := getTestDB(t)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	assert.NoError(t, err)
	accountStore := NewAccountStore(db)
	transactionStore := NewTransactionStore(db)

	acc, err := accountStore.Create(accountFixture1)
	assert.NoError(t, err)
	transactionFixture1.AccountID = acc.AccountID
	_, err = transactionStore.Create(transactionFixture1)
	assert.NoError(t, err)

	retrievedTransaction, err := transactionStore.Retrieve(transactionFixture1.TransactionID)
	assert.NoError(t, err)
	assert.Equal(t, transactionFixture1, retrievedTransaction)
}
