package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/NavaneethWKT/CapStone_GO_Lang/server/internal/repository"
)

type WalletService struct {
	walletRepo *repository.WalletRepository
}

// NewWalletService creates a new wallet service
func NewWalletService(walletRepo *repository.WalletRepository) *WalletService {
	return &WalletService{
		walletRepo: walletRepo,
	}
}

// GetBalance retrieves the current balance for a user
func (s *WalletService) GetBalance(userID int) (float64, error) {
	if userID <= 0 {
		return 0, errors.New("invalid user ID")
	}

	balance, err := s.walletRepo.GetBalance(userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get wallet balance: %w", err)
	}

	return balance, nil
}

// ValidateSufficientBalance checks if user has enough balance for a transaction
func (s *WalletService) ValidateSufficientBalance(userID int, amount float64) error {
	if userID <= 0 {
		return errors.New("invalid user ID")
	}

	if amount <= 0 {
		return errors.New("invalid amount")
	}

	balance, err := s.walletRepo.GetBalance(userID)
	if err != nil {
		return fmt.Errorf("failed to get wallet balance: %w", err)
	}

	if balance < amount {
		return errors.New("insufficient wallet balance")
	}

	return nil
}

// DeductBalance deducts amount from user's wallet (used in transactions)
func (s *WalletService) DeductBalance(tx *sql.Tx, userID int, amount float64) error {
	if userID <= 0 {
		return errors.New("invalid user ID")
	}

	if amount <= 0 {
		return errors.New("invalid amount")
	}

	// Validate balance before deducting
	if err := s.ValidateSufficientBalance(userID, amount); err != nil {
		return err
	}

	// Deduct the amount (negative value)
	err := s.walletRepo.UpdateBalance(tx, userID, -amount)
	if err != nil {
		return fmt.Errorf("failed to deduct wallet balance: %w", err)
	}

	return nil
}

// AddBalance adds amount to user's wallet (for top-ups, refunds)
func (s *WalletService) AddBalance(tx *sql.Tx, userID int, amount float64) error {
	if userID <= 0 {
		return errors.New("invalid user ID")
	}

	if amount <= 0 {
		return errors.New("invalid amount")
	}

	// Add the amount (positive value)
	err := s.walletRepo.UpdateBalance(tx, userID, amount)
	if err != nil {
		return fmt.Errorf("failed to add wallet balance: %w", err)
	}

	return nil
}

