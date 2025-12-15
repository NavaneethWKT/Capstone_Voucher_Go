package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/NavaneethWKT/CapStone_GO_Lang/server/internal/model"
	"github.com/NavaneethWKT/CapStone_GO_Lang/server/internal/repository"
)

type PaymentService struct {
	db                *sql.DB
	userService       *UserService
	voucherService    *VoucherService
	walletService     *WalletService
	voucherRepo       *repository.VoucherRepository
	transactionRepo   *repository.TransactionRepository
	mockUPI           *MockUPI
}

// NewPaymentService creates a new payment service
func NewPaymentService(
	db *sql.DB,
	userService *UserService,
	voucherService *VoucherService,
	walletService *WalletService,
	voucherRepo *repository.VoucherRepository,
	transactionRepo *repository.TransactionRepository,
	mockUPI *MockUPI,
) *PaymentService {
	return &PaymentService{
		db:              db,
		userService:     userService,
		voucherService:  voucherService,
		walletService:   walletService,
		voucherRepo:     voucherRepo,
		transactionRepo: transactionRepo,
		mockUPI:         mockUPI,
	}
}

// BuyVoucher orchestrates the complete voucher purchase flow with ACID transactions
func (s *PaymentService) BuyVoucher(userID, voucherID int) (*model.Transaction, error) {
	// Step 1: Validate user exists
	if err := s.userService.ValidateUserExists(userID); err != nil {
		return nil, err
	}

	// Step 2: Validate voucher is available
	if err := s.voucherService.ValidateVoucherAvailable(voucherID); err != nil {
		return nil, err
	}

	// Step 3: Get voucher details (need price)
	voucher, err := s.voucherService.GetVoucherByID(voucherID)
	if err != nil {
		return nil, err
	}

	amount := voucher.Price

	// Step 4: Validate sufficient balance
	if err := s.walletService.ValidateSufficientBalance(userID, amount); err != nil {
		return nil, err
	}

	// Step 5: Start database transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	// Defer rollback in case of error
	defer tx.Rollback()

	// Step 6: Deduct from wallet
	if err := s.walletService.DeductBalance(tx, userID, amount); err != nil {
		return nil, err
	}

	// Step 7: Create transaction record (pending status)
	transaction := &model.Transaction{
		UserID:           userID,
		VoucherID:        &voucherID,
		Amount:           amount,
		TransactionType:  model.TransactionTypePurchase,
		PaymentStatus:    model.PaymentStatusPending,
		PaymentTxnID:     nil,
	}

	if err := s.transactionRepo.CreateTransaction(tx, transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}

	// Step 8: Update voucher quantity (decrease by 1)
	if err := s.voucherRepo.UpdateVoucherQuantity(tx, voucherID, -1); err != nil {
		return nil, fmt.Errorf("failed to update voucher quantity: %w", err)
	}

	// Step 9: Process payment via Mock UPI
	paymentResult, err := s.mockUPI.ProcessPayment(amount, userID, transaction.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to process payment: %w", err)
	}

	// Step 10: Update transaction status based on payment result
	var paymentTxnID *string
	if paymentResult.Success {
		paymentTxnID = &paymentResult.PaymentTxnID
		if err := s.transactionRepo.UpdateTransactionStatus(tx, transaction.ID, model.PaymentStatusSuccess, paymentTxnID); err != nil {
			return nil, fmt.Errorf("failed to update transaction status: %w", err)
		}
		transaction.PaymentStatus = model.PaymentStatusSuccess
		transaction.PaymentTxnID = paymentTxnID
	} else {
		// Payment failed - rollback will happen automatically
		if err := s.transactionRepo.UpdateTransactionStatus(tx, transaction.ID, model.PaymentStatusFailed, nil); err != nil {
			return nil, fmt.Errorf("failed to update transaction status: %w", err)
		}
		transaction.PaymentStatus = model.PaymentStatusFailed
		return nil, errors.New("payment processing failed")
	}

	// Step 11: Commit transaction (all operations succeeded)
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return transaction, nil
}

