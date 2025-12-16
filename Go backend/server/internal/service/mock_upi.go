package service

import (
	"fmt"
	"math/rand"
	"time"
)

// PaymentResult represents the result of a payment processing
type PaymentResult struct {
	Success      bool
	PaymentTxnID string
	Message      string
}

// MockUPI simulates a payment gateway
type MockUPI struct {
	successRate float64 // Success rate (0.0 to 1.0), default 0.95 (95%)
}

// NewMockUPI creates a new Mock UPI instance
func NewMockUPI(successRate float64) *MockUPI {
	if successRate <= 0 || successRate > 1 {
		successRate = 0.95 // Default 95% success rate
	}
	return &MockUPI{
		successRate: successRate,
	}
}

// ProcessPayment simulates payment processing with network delay
func (m *MockUPI) ProcessPayment(amount float64, userID, transactionID int) (*PaymentResult, error) {
	// Create random generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	// Simulate network delay (100-500ms)
	delay := time.Duration(100+r.Intn(400)) * time.Millisecond
	time.Sleep(delay)

	// Validate amount
	if amount <= 0 {
		return &PaymentResult{
			Success: false,
			Message: "invalid payment amount",
		}, nil
	}

	// Simulate success/failure based on success rate
	success := r.Float64() < m.successRate

	if success {
		// Generate mock payment transaction ID
		paymentTxnID := fmt.Sprintf("UPI_%d_%d_%d", userID, transactionID, time.Now().Unix())
		
		return &PaymentResult{
			Success:      true,
			PaymentTxnID: paymentTxnID,
			Message:      "payment processed successfully",
		}, nil
	}

	// Payment failed
	return &PaymentResult{
		Success: false,
		Message: "payment processing failed",
	}, nil
}

