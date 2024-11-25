package stores

import (
	"database/sql"
	"fmt"
	"transactions/models"

	sq "github.com/Masterminds/squirrel"
)

type TransactionStore struct {
	db          *sql.DB
	insertQuery sq.InsertBuilder
	selectQuery sq.SelectBuilder
	updateQuery sq.UpdateBuilder
}

func NewTransactionStore(db *sql.DB) *TransactionStore {
	return &TransactionStore{
		db: db,
		insertQuery: psql.Insert("transaction").
			Columns("account_id", "amount", "operation_type", "balance").
			Suffix("RETURNING \"id\", \"created_at\""),
		selectQuery: psql.Select("id", "account_id", "amount", "operation_type, created_at, balance").
			From("transaction"),
		updateQuery: psql.Update("balance").
			From("transaction"),
	}
}

func (s *TransactionStore) Create(transaction *models.Transaction) (*models.Transaction, error) {
	transaction.Balance = transaction.Amount
	if transaction.OperationTypeID != 1 {
		// calc balance
		s.getPaymentBalance(transaction)
	}

	err := s.insertQuery.Values(
		transaction.AccountID,
		transaction.Amount,
		transaction.OperationTypeID,
		transaction.Balance,
	).RunWith(s.db).QueryRow().Scan(&transaction.TransactionID, &transaction.CreatedAt)

	return transaction, err
}

func (s *TransactionStore) Retrieve(id int64) (*models.Transaction, error) {
	row := s.selectQuery.Where(sq.Eq{"id": id}).RunWith(s.db).QueryRow()
	transaction := &models.Transaction{}
	err := row.Scan(&transaction.TransactionID, &transaction.AccountID, &transaction.Amount, &transaction.OperationTypeID, &transaction.CreatedAt, &transaction.Balance)
	return transaction, err
}

func (s *TransactionStore) Search(account_id int64) ([]models.Transaction, error) {

	q := s.selectQuery.Where(sq.Eq{"account_id": account_id}, sq.Lt{"amount": 0})
	sql, params, err := q.ToSql()
	fmt.Println("sql-", sql, " params-", params)

	rows, err := s.db.Query(sql, params...)
	result := []models.Transaction{}
	fmt.Println("err=", err)
	for rows.Next() {
		t := models.Transaction{}
		err := rows.Scan(t)

		if err != nil {
			return result, err
		}
		result = append(result, t)
	}
	return result, err
}

func (s *TransactionStore) Update(trans models.Transaction) (models.Transaction, error) {
	_, err := s.updateQuery.Where(sq.Eq{"id": trans.TransactionID}).RunWith(s.db).Exec()

	return trans, err
}

func (s *TransactionStore) getPaymentBalance(transaction *models.Transaction) (*models.Transaction, error) {
	undischarged, err := s.Search(transaction.AccountID)
	if err != nil {
		return transaction, err
	}

	// iterate through transactions that need discharge, sub from trans
	for _, undischarged_item := range undischarged {
		// if we have plenty of mondy
		if transaction.Balance >= undischarged_item.Balance {
			transaction.Balance = transaction.Balance - undischarged_item.Balance
		} else { //not engough to discharge
			undischarged_item.Balance = transaction.Balance - undischarged_item.Balance
			transaction.Balance = 0
		}
		s.Update(undischarged_item)
	}
	return transaction, err
}
