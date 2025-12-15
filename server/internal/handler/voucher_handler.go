package handler

import (
	"context"
	"time"

	"github.com/NavaneethWKT/CapStone_GO_Lang/protoc"
	"github.com/NavaneethWKT/CapStone_GO_Lang/server/internal/errors"
	"github.com/NavaneethWKT/CapStone_GO_Lang/server/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type VoucherHandler struct {
	voucherService *service.VoucherService
}

// NewVoucherHandler creates a new voucher handler
func NewVoucherHandler(voucherService *service.VoucherService) *VoucherHandler {
	return &VoucherHandler{
		voucherService: voucherService,
	}
}

// Search searches for vouchers with optional filters
func (h *VoucherHandler) Search(ctx context.Context, req *protoc.SearchRequest) (*protoc.SearchResponse, error) {
	// Extract filters from request
	category := req.GetCategory()
	var minPrice, maxPrice *float64

	if req.GetMinPrice() > 0 {
		price := req.GetMinPrice()
		minPrice = &price
	}

	if req.GetMaxPrice() > 0 {
		price := req.GetMaxPrice()
		maxPrice = &price
	}

	// Call service
	vouchers, err := h.voucherService.SearchVouchers(category, minPrice, maxPrice)
	if err != nil {
		return nil, h.handleError(err)
	}

	// Convert domain models to gRPC messages
	pbVouchers := make([]*protoc.Voucher, 0, len(vouchers))
	for _, v := range vouchers {
		pbVouchers = append(pbVouchers, &protoc.Voucher{
			Id:          int32(v.ID),
			Name:        v.Name,
			Description: v.Description,
			Category:    v.Category,
			Price:       v.Price,
			Quantity:    int32(v.Quantity),
			ValidFrom:   v.ValidFrom.Format(time.RFC3339),
			ValidTo:     v.ValidTo.Format(time.RFC3339),
			CreatedAt:   v.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   v.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &protoc.SearchResponse{
		Vouchers: pbVouchers,
	}, nil
}

// handleError converts application errors to gRPC status errors
func (h *VoucherHandler) handleError(err error) error {
	switch err {
	case errors.ErrVoucherNotFound:
		return status.Error(codes.NotFound, err.Error())
	case errors.ErrVoucherOutOfStock:
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.ErrVoucherExpired:
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.ErrInvalidVoucherID:
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.ErrInvalidPrice:
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.ErrUserNotFound:
		return status.Error(codes.NotFound, err.Error())
	case errors.ErrInvalidUserID:
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.ErrInsufficientBalance:
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.ErrPaymentFailed:
		return status.Error(codes.Internal, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}

