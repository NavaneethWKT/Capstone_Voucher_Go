package errors

import (
	"errors"
	"fmt"
)

var (
	// User errors
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidUserID     = errors.New("invalid user ID")

	// Voucher errors
	ErrVoucherNotFound      = errors.New("voucher not found")
	ErrVoucherOutOfStock    = errors.New("voucher out of stock")
	ErrVoucherExpired       = errors.New("voucher expired")
	ErrVoucherNotAvailable  = errors.New("voucher not available")
	ErrInvalidVoucherID      = errors.New("invalid voucher ID")

	// Wallet errors
	ErrWalletNotFound       = errors.New("wallet not found")
	ErrInsufficientBalance  = errors.New("insufficient wallet balance")
	ErrInvalidAmount        = errors.New("invalid amount")
	ErrNegativeBalance      = errors.New("balance cannot be negative")

	// Transaction errors
	ErrTransactionNotFound  = errors.New("transaction not found")
	ErrTransactionFailed    = errors.New("transaction failed")
	ErrInvalidTransactionID = errors.New("invalid transaction ID")
	ErrPaymentFailed        = errors.New("payment processing failed")

	// Database errors
	ErrDatabaseConnection = errors.New("database connection error")
	ErrDatabaseQuery      = errors.New("database query error")
	ErrDatabaseTransaction = errors.New("database transaction error")

	// Validation errors
	ErrInvalidInput      = errors.New("invalid input")
	ErrMissingRequiredField = errors.New("missing required field")
	ErrInvalidEmail      = errors.New("invalid email format")
	ErrInvalidPrice      = errors.New("invalid price")
	ErrInvalidQuantity   = errors.New("invalid quantity")
)

// Error wrapper for adding context
type AppError struct {
	Err     error
	Message string
	Code    string
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Helper functions to create errors with context
func NewAppError(err error, message string, code string) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}

func WrapError(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}

// Error codes for API responses
const (
	CodeUserNotFound      = "USER_NOT_FOUND"
	CodeUserAlreadyExists = "USER_ALREADY_EXISTS"
	CodeInvalidUserID     = "INVALID_USER_ID"

	CodeVoucherNotFound     = "VOUCHER_NOT_FOUND"
	CodeVoucherOutOfStock   = "VOUCHER_OUT_OF_STOCK"
	CodeVoucherExpired      = "VOUCHER_EXPIRED"
	CodeVoucherNotAvailable = "VOUCHER_NOT_AVAILABLE"
	CodeInvalidVoucherID    = "INVALID_VOUCHER_ID"

	CodeWalletNotFound      = "WALLET_NOT_FOUND"
	CodeInsufficientBalance = "INSUFFICIENT_BALANCE"
	CodeInvalidAmount       = "INVALID_AMOUNT"
	CodeNegativeBalance     = "NEGATIVE_BALANCE"

	CodeTransactionNotFound  = "TRANSACTION_NOT_FOUND"
	CodeTransactionFailed    = "TRANSACTION_FAILED"
	CodeInvalidTransactionID = "INVALID_TRANSACTION_ID"
	CodePaymentFailed        = "PAYMENT_FAILED"

	CodeDatabaseConnection  = "DATABASE_CONNECTION_ERROR"
	CodeDatabaseQuery       = "DATABASE_QUERY_ERROR"
	CodeDatabaseTransaction = "DATABASE_TRANSACTION_ERROR"

	CodeInvalidInput        = "INVALID_INPUT"
	CodeMissingRequiredField = "MISSING_REQUIRED_FIELD"
	CodeInvalidEmail        = "INVALID_EMAIL"
	CodeInvalidPrice        = "INVALID_PRICE"
	CodeInvalidQuantity     = "INVALID_QUANTITY"
)

