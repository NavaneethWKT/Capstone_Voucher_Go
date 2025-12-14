package repository

import (
	"database/sql"
	"time"

	"github.com/NavaneethWKT/CapStone_GO_Lang/server/internal/model"
)

type TransactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// CreateTransaction creates a new transaction
func (r *TransactionRepository) CreateTransaction(tx *sql.Tx, transaction *model.Transaction) error {
	query := `
		INSERT INTO transactions (user_id, voucher_id, amount, transaction_type, payment_status, payment_txn_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	now := time.Now()
	transaction.CreatedAt = now
	transaction.UpdatedAt = now

	var err error
	if tx != nil {
		err = tx.QueryRow(
			query,
			transaction.UserID,
			transaction.VoucherID,
			transaction.Amount,
			transaction.TransactionType,
			transaction.PaymentStatus,
			transaction.PaymentTxnID,
			now,
			now,
		).Scan(&transaction.ID)
	} else {
		err = r.db.QueryRow(
			query,
			transaction.UserID,
			transaction.VoucherID,
			transaction.Amount,
			transaction.TransactionType,
			transaction.PaymentStatus,
			transaction.PaymentTxnID,
			now,
			now,
		).Scan(&transaction.ID)
	}

	return err
}

// GetTransactionsByUserID retrieves all transactions for a user with optional filters
func (r *TransactionRepository) GetTransactionsByUserID(userID int, limit, offset int) ([]*model.Transaction, error) {
	query := `
		SELECT id, user_id, voucher_id, amount, transaction_type, payment_status, payment_txn_id, created_at, updated_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		transaction := &model.Transaction{}
		var voucherID sql.NullInt64
		var paymentTxnID sql.NullString

		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&voucherID,
			&transaction.Amount,
			&transaction.TransactionType,
			&transaction.PaymentStatus,
			&paymentTxnID,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if voucherID.Valid {
			vID := int(voucherID.Int64)
			transaction.VoucherID = &vID
		}

		if paymentTxnID.Valid {
			transaction.PaymentTxnID = &paymentTxnID.String
		}

		transactions = append(transactions, transaction)
	}

	return transactions, rows.Err()
}

// UpdateTransactionStatus updates the payment status and payment transaction ID
func (r *TransactionRepository) UpdateTransactionStatus(tx *sql.Tx, transactionID int, status model.PaymentStatus, paymentTxnID *string) error {
	query := `
		UPDATE transactions
		SET payment_status = $1, payment_txn_id = $2, updated_at = $3
		WHERE id = $4
	`

	var err error
	if tx != nil {
		_, err = tx.Exec(query, status, paymentTxnID, time.Now(), transactionID)
	} else {
		_, err = r.db.Exec(query, status, paymentTxnID, time.Now(), transactionID)
	}

	return err
}

// GetTransactionByID retrieves a transaction by ID
func (r *TransactionRepository) GetTransactionByID(id int) (*model.Transaction, error) {
	query := `
		SELECT id, user_id, voucher_id, amount, transaction_type, payment_status, payment_txn_id, created_at, updated_at
		FROM transactions
		WHERE id = $1
	`

	transaction := &model.Transaction{}
	var voucherID sql.NullInt64
	var paymentTxnID sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&transaction.ID,
		&transaction.UserID,
		&voucherID,
		&transaction.Amount,
		&transaction.TransactionType,
		&transaction.PaymentStatus,
		&paymentTxnID,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Transaction not found
		}
		return nil, err
	}

	if voucherID.Valid {
		vID := int(voucherID.Int64)
		transaction.VoucherID = &vID
	}

	if paymentTxnID.Valid {
		transaction.PaymentTxnID = &paymentTxnID.String
	}

	return transaction, nil
}

