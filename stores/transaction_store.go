package stores

import (
	"database/sql"
	"transactions/models"

	sq "github.com/Masterminds/squirrel"
)

type TransactionStore struct {
	db          *sql.DB
	insertQuery sq.InsertBuilder
	selectQuery sq.SelectBuilder
}

func NewTransactionStore(db *sql.DB) *TransactionStore {
	return &TransactionStore{
		db: db,
		insertQuery: psql.Insert("transaction").
			Columns("account_id", "amount", "operation_type").
			Suffix("RETURNING \"id\", \"created_at\""),
		selectQuery: psql.Select("id", "account_id", "amount", "operation_type, created_at").
			From("transaction"),
	}
}

func (s *TransactionStore) Create(transaction *models.Transaction) (*models.Transaction, error) {
	err := s.insertQuery.Values(
		transaction.AccountID,
		transaction.Amount,
		transaction.OperationTypeID,
	).RunWith(s.db).QueryRow().Scan(&transaction.TransactionID, &transaction.CreatedAt)

	return transaction, err
}

func (s *TransactionStore) Retrieve(id int64) (*models.Transaction, error) {
	row := s.selectQuery.Where(sq.Eq{"id": id}).RunWith(s.db).QueryRow()
	transaction := &models.Transaction{}
	err := row.Scan(&transaction.TransactionID, &transaction.AccountID, &transaction.Amount, &transaction.OperationTypeID, &transaction.CreatedAt)
	return transaction, err
}
