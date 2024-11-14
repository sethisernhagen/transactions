package stores

import (
	"database/sql"
	"testing"
	"transactions/models"
)

var (
	accountFixture1 = &models.Account{
		AccountID:      -1,
		DocumentNumber: "12345678900",
	}
	transactionFixture1 = &models.Transaction{
		TransactionID:   -1,
		AccountID:       -1,
		OperationTypeID: 1,
		Amount:          100.0,
	}
)

func getTestDB(t *testing.T) (*sql.DB, error) {
	db, err := GetDB("localhost", 5432, "postgres", "example", "postgres")
	if err != nil {
		t.Fatal(err.Error())
	}
	truncateTables(t, db)
	return db, err
}

func truncateTables(t *testing.T, db *sql.DB) {
	_, err := db.Exec("DELETE FROM transaction")
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = db.Exec("DELETE FROM account")
	if err != nil {
		t.Fatal(err.Error())
	}
}
