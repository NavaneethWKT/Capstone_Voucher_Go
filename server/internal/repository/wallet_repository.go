package repository

import (
	"database/sql"
	"time"

	"github.com/NavaneethWKT/CapStone_GO_Lang/server/internal/model"
)

type WalletRepository struct {
	db *sql.DB
}

// NewWalletRepository creates a new wallet repository
func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

// GetWalletByUserID retrieves wallet for a user
func (r *WalletRepository) GetWalletByUserID(userID int) (*model.Wallet, error) {
	query := `
		SELECT id, user_id, balance, created_at, updated_at
		FROM wallets
		WHERE user_id = $1
	`

	wallet := &model.Wallet{}
	err := r.db.QueryRow(query, userID).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Balance,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Wallet not found
		}
		return nil, err
	}

	return wallet, nil
}

// CreateWallet creates a new wallet for a user
func (r *WalletRepository) CreateWallet(userID int, initialBalance float64) (*model.Wallet, error) {
	query := `
		INSERT INTO wallets (user_id, balance, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	now := time.Now()
	wallet := &model.Wallet{
		UserID:    userID,
		Balance:   initialBalance,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := r.db.QueryRow(query, userID, initialBalance, now, now).Scan(&wallet.ID)
	if err != nil {
		return nil, err
	}

	return wallet, nil
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

