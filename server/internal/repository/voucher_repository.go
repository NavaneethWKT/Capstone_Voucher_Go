package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/NavaneethWKT/CapStone_GO_Lang/server/internal/model"
)

type VoucherRepository struct {
	db *sql.DB
}

// NewVoucherRepository creates a new voucher repository
func NewVoucherRepository(db *sql.DB) *VoucherRepository {
	return &VoucherRepository{db: db}
}

// SearchVouchers searches for vouchers with optional filters
func (r *VoucherRepository) SearchVouchers(category string, minPrice, maxPrice *float64, limit, offset int) ([]*model.Voucher, error) {
	query := `
		SELECT id, name, description, category, price, quantity, valid_from, valid_to, created_at, updated_at
		FROM vouchers
		WHERE valid_from <= $1 AND valid_to >= $1 AND quantity > 0
	`
	args := []interface{}{time.Now()}
	argPos := 2

	if category != "" {
		query += fmt.Sprintf(" AND category = $%d", argPos)
		args = append(args, category)
		argPos++
	}

	if minPrice != nil {
		query += fmt.Sprintf(" AND price >= $%d", argPos)
		args = append(args, *minPrice)
		argPos++
	}

	if maxPrice != nil {
		query += fmt.Sprintf(" AND price <= $%d", argPos)
		args = append(args, *maxPrice)
		argPos++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vouchers []*model.Voucher
	for rows.Next() {
		voucher := &model.Voucher{}
		err := rows.Scan(
			&voucher.ID,
			&voucher.Name,
			&voucher.Description,
			&voucher.Category,
			&voucher.Price,
			&voucher.Quantity,
			&voucher.ValidFrom,
			&voucher.ValidTo,
			&voucher.CreatedAt,
			&voucher.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		vouchers = append(vouchers, voucher)
	}

	return vouchers, rows.Err()
}

// GetVoucherByID retrieves a voucher by ID
func (r *VoucherRepository) GetVoucherByID(id int) (*model.Voucher, error) {
	query := `
		SELECT id, name, description, category, price, quantity, valid_from, valid_to, created_at, updated_at
		FROM vouchers
		WHERE id = $1
	`

	voucher := &model.Voucher{}
	err := r.db.QueryRow(query, id).Scan(
		&voucher.ID,
		&voucher.Name,
		&voucher.Description,
		&voucher.Category,
		&voucher.Price,
		&voucher.Quantity,
		&voucher.ValidFrom,
		&voucher.ValidTo,
		&voucher.CreatedAt,
		&voucher.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Voucher not found
		}
		return nil, err
	}

	return voucher, nil
}

// UpdateVoucherQuantity updates the quantity of a voucher (used in transactions)
func (r *VoucherRepository) UpdateVoucherQuantity(tx *sql.Tx, voucherID int, quantityChange int) error {
	query := `
		UPDATE vouchers
		SET quantity = quantity + $1, updated_at = $2
		WHERE id = $3
	`

	var err error
	if tx != nil {
		_, err = tx.Exec(query, quantityChange, time.Now(), voucherID)
	} else {
		_, err = r.db.Exec(query, quantityChange, time.Now(), voucherID)
	}

	return err
}

