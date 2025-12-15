package repository

import (
	"database/sql"
	"time"
)

type WalletRepository struct {
	db *sql.DB
}

// NewWalletRepository creates a new wallet repository
func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

// UpdateBalance updates wallet balance (supports transactions for ACID)
func (r *WalletRepository) UpdateBalance(tx *sql.Tx, userID int, amountChange float64) error {
	query := `
		UPDATE wallets
		SET balance = balance + $1, updated_at = $2
		WHERE user_id = $3
	`

	var err error
	if tx != nil {
		_, err = tx.Exec(query, amountChange, time.Now(), userID)
	} else {
		_, err = r.db.Exec(query, amountChange, time.Now(), userID)
	}

	return err
}

// GetBalance retrieves current balance for a user
func (r *WalletRepository) GetBalance(userID int) (float64, error) {
	query := `
		SELECT balance
		FROM wallets
		WHERE user_id = $1
	`

	var balance float64
	err := r.db.QueryRow(query, userID).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // Wallet doesn't exist, return 0
		}
		return 0, err
	}

	return balance, nil
}

