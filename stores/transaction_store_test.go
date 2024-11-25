package stores

import (
	"testing"
	"transactions/models"

	"github.com/stretchr/testify/assert"
)

func TestTransactionStore_CreatePurchase(t *testing.T) {
	db, err := GetTestDB(t)
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

func TestTransactionStore_CreateCredit(t *testing.T) {
	db, err := GetTestDB(t)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	assert.NoError(t, err)
	accountStore := NewAccountStore(db)
	transactionStore := NewTransactionStore(db)

	acc, err := accountStore.Create(accountFixture1)
	assert.NoError(t, err)
	transactionCredit1.AccountID = acc.AccountID
	createdTransaction, err := transactionStore.Create(transactionCredit1)
	assert.NoError(t, err)
	assert.Equal(t, transactionCredit1, createdTransaction)
}

func TestTransactionStore_CreateCreditWithDischarge(t *testing.T) {
	db, err := GetTestDB(t)
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
	purchaseTransaction, err := transactionStore.Create(transactionFixture1)
	assert.NoError(t, err)
	transactionCredit1.AccountID = acc.AccountID
	result, err := transactionStore.Create(transactionCredit1)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), result.Balance)
	// check balance of purchase transaction
	purchaseResult, err := transactionStore.Retrieve(purchaseTransaction.TransactionID)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), purchaseResult.Balance)
}

func TestTransactionStore_CreateCreditWithDischargeWithLeftoverBalance(t *testing.T) {
	db, err := GetTestDB(t)
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
	purchaseTransaction, err := transactionStore.Create(transactionFixture1)
	assert.NoError(t, err)
	transactionCredit1.AccountID = acc.AccountID
	result, err := transactionStore.Create(
		&models.Transaction{
			AccountID:       acc.AccountID,
			OperationTypeID: 4,
			Amount:          50.0,
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), result.Balance)
	// check balance of purchase transaction
	purchaseResult, err := transactionStore.Retrieve(purchaseTransaction.TransactionID)
	assert.NoError(t, err)
	assert.Equal(t, float64(-50), purchaseResult.Balance)
}

func TestTransactionStore_CreateCreditWithPartialDischarge(t *testing.T) {
	db, err := GetTestDB(t)
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
	purchaseTransaction, err := transactionStore.Create(transactionFixture1)
	assert.NoError(t, err)
	transactionCredit1.AccountID = acc.AccountID
	result, err := transactionStore.Create(
		&models.Transaction{
			AccountID:       acc.AccountID,
			OperationTypeID: 4,
			Amount:          50.00,
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), result.Balance)
	// check balance of purchase transaction
	purchaseResult, err := transactionStore.Retrieve(purchaseTransaction.TransactionID)
	assert.NoError(t, err)
	assert.Equal(t, float64(-50), purchaseResult.Balance)
}

func TestTransactionStore_Example2(t *testing.T) {
	db, err := GetTestDB(t)
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
	trans1, err := transactionStore.Create(
		&models.Transaction{
			AccountID:       acc.AccountID,
			OperationTypeID: 1,
			Amount:          -50.00,
		},
	)
	assert.NoError(t, err)
	trans2, err := transactionStore.Create(
		&models.Transaction{
			AccountID:       acc.AccountID,
			OperationTypeID: 1,
			Amount:          -23.50,
		},
	)
	assert.NoError(t, err)
	trans3, err := transactionStore.Create(
		&models.Transaction{
			AccountID:       acc.AccountID,
			OperationTypeID: 1,
			Amount:          -18.70,
		},
	)
	assert.NoError(t, err)

	result, err := transactionStore.Create(
		&models.Transaction{
			AccountID:       acc.AccountID,
			OperationTypeID: 4,
			Amount:          60.00,
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), result.Balance)
	// check balance of purchase transactions
	trans1Result, err := transactionStore.Retrieve(trans1.TransactionID)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), trans1Result.Balance)
	trans2Result, err := transactionStore.Retrieve(trans2.TransactionID)
	assert.NoError(t, err)
	assert.Equal(t, float64(-13.50), trans2Result.Balance)
	trans3Result, err := transactionStore.Retrieve(trans3.TransactionID)
	assert.NoError(t, err)
	assert.Equal(t, float64(-18.70), trans3Result.Balance)
	//check the balance of credit transaction
	creditResult, err := transactionStore.Retrieve(result.TransactionID)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), creditResult.Balance)
}

func TestTransactionStore_Example3(t *testing.T) {
	db, err := GetTestDB(t)
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
	trans1, err := transactionStore.Create(
		&models.Transaction{
			AccountID:       acc.AccountID,
			OperationTypeID: 1,
			Amount:          -50.00,
		},
	)
	assert.NoError(t, err)
	trans2, err := transactionStore.Create(
		&models.Transaction{
			AccountID:       acc.AccountID,
			OperationTypeID: 1,
			Amount:          -23.50,
		},
	)
	assert.NoError(t, err)
	trans3, err := transactionStore.Create(
		&models.Transaction{
			AccountID:       acc.AccountID,
			OperationTypeID: 1,
			Amount:          -18.70,
		},
	)
	assert.NoError(t, err)
	trans4, err := transactionStore.Create(
		&models.Transaction{
			AccountID:       acc.AccountID,
			OperationTypeID: 4,
			Amount:          60.00,
		},
	)
	assert.NoError(t, err)

	result, err := transactionStore.Create(
		&models.Transaction{
			AccountID:       acc.AccountID,
			OperationTypeID: 4,
			Amount:          100.00,
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, float64(67.80), result.Balance)

	// check balance of purchase transactions
	trans1Result, err := transactionStore.Retrieve(trans1.TransactionID)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), trans1Result.Balance)
	trans2Result, err := transactionStore.Retrieve(trans2.TransactionID)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), trans2Result.Balance)
	trans3Result, err := transactionStore.Retrieve(trans3.TransactionID)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), trans3Result.Balance)
	trans4Result, err := transactionStore.Retrieve(trans4.TransactionID)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), trans4Result.Balance)
	//check the balance of credit transaction
	creditResult, err := transactionStore.Retrieve(result.TransactionID)
	assert.NoError(t, err)
	assert.Equal(t, float64(67.80), creditResult.Balance)

}

func TestTransactionStore_Retrieve(t *testing.T) {
	db, err := GetTestDB(t)
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
