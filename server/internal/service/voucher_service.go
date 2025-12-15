package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/NavaneethWKT/CapStone_GO_Lang/server/internal/model"
	"github.com/NavaneethWKT/CapStone_GO_Lang/server/internal/repository"
)

type VoucherService struct {
	voucherRepo *repository.VoucherRepository
}

// NewVoucherService creates a new voucher service
func NewVoucherService(voucherRepo *repository.VoucherRepository) *VoucherService {
	return &VoucherService{
		voucherRepo: voucherRepo,
	}
}

// SearchVouchers searches for vouchers with optional filters
func (s *VoucherService) SearchVouchers(category string, minPrice, maxPrice *float64) ([]*model.Voucher, error) {
	// Validate price range if provided
	if minPrice != nil && *minPrice < 0 {
		return nil, errors.New("invalid price")
	}
	if maxPrice != nil && *maxPrice < 0 {
		return nil, errors.New("invalid price")
	}
	if minPrice != nil && maxPrice != nil && *minPrice > *maxPrice {
		return nil, errors.New("invalid price")
	}

	vouchers, err := s.voucherRepo.SearchVouchers(category, minPrice, maxPrice)
	if err != nil {
		return nil, fmt.Errorf("failed to search vouchers: %w", err)
	}

	return vouchers, nil
}

// GetVoucherByID retrieves a voucher by ID
func (s *VoucherService) GetVoucherByID(voucherID int) (*model.Voucher, error) {
	if voucherID <= 0 {
		return nil, errors.New("invalid voucher ID")
	}

	voucher, err := s.voucherRepo.GetVoucherByID(voucherID)
	if err != nil {
		return nil, fmt.Errorf("failed to get voucher: %w", err)
	}

	if voucher == nil {
		return nil, errors.New("voucher not found")
	}

	return voucher, nil
}

// ValidateVoucherAvailable checks if voucher exists, is in stock, and not expired
func (s *VoucherService) ValidateVoucherAvailable(voucherID int) error {
	if voucherID <= 0 {
		return errors.New("invalid voucher ID")
	}

	voucher, err := s.voucherRepo.GetVoucherByID(voucherID)
	if err != nil {
		return fmt.Errorf("failed to get voucher: %w", err)
	}

	if voucher == nil {
		return errors.New("voucher not found")
	}

	// Check if voucher is expired
	now := time.Now()
	if now.Before(voucher.ValidFrom) || now.After(voucher.ValidTo) {
		return errors.New("voucher expired")
	}

	// Check if voucher is out of stock
	if voucher.Quantity <= 0 {
		return errors.New("voucher out of stock")
	}

	return nil
}

