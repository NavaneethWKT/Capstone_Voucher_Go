package model

import "time"

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypePurchase TransactionType = "purchase"
	TransactionTypeRefund   TransactionType = "refund"
	TransactionTypeTopUp    TransactionType = "topup"
)

// PaymentStatus represents the payment status
type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusSuccess PaymentStatus = "success"
	PaymentStatusFailed  PaymentStatus = "failed"
)

// Transaction represents a transaction in the system
type Transaction struct {
	ID                int             `json:"id" db:"id"`
	UserID            int             `json:"user_id" db:"user_id"`
	VoucherID         *int            `json:"voucher_id,omitempty" db:"voucher_id"` // Nullable
	Amount            float64         `json:"amount" db:"amount"`
	TransactionType  TransactionType `json:"transaction_type" db:"transaction_type"`
	PaymentStatus     PaymentStatus    `json:"payment_status" db:"payment_status"`
	PaymentTxnID      *string          `json:"payment_txn_id,omitempty" db:"payment_txn_id"` // Nullable, from Mock UPI
	CreatedAt         time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at" db:"updated_at"`
}

