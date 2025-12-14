package model

import "time"

// Voucher represents a voucher in the system
type Voucher struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Category    string    `json:"category" db:"category"`
	Price       float64   `json:"price" db:"price"`
	Quantity    int       `json:"quantity" db:"quantity"`
	ValidFrom   time.Time `json:"valid_from" db:"valid_from"`
	ValidTo     time.Time `json:"valid_to" db:"valid_to"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

