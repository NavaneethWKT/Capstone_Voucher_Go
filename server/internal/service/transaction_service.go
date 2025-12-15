package service

import (
	"github.com/NavaneethWKT/CapStone_GO_Lang/server/internal/errors"
	"github.com/NavaneethWKT/CapStone_GO_Lang/server/internal/model"
	"github.com/NavaneethWKT/CapStone_GO_Lang/server/internal/repository"
)

type TransactionService struct {
	transactionRepo *repository.TransactionRepository
	userService     *UserService
}

// NewTransactionService creates a new transaction service
func NewTransactionService(transactionRepo *repository.TransactionRepository, userService *UserService) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		userService:     userService,
	}
}

// ListTransactions retrieves all transactions for a user
func (s *TransactionService) ListTransactions(userID int) ([]*model.Transaction, error) {
	// Validate user exists
	if err := s.userService.ValidateUserExists(userID); err != nil {
		return nil, err
	}

	transactions, err := s.transactionRepo.GetTransactionsByUserID(userID)
	if err != nil {
		return nil, errors.WrapError(err, "failed to get transactions")
	}

	return transactions, nil
}

