package handler

import (
	"context"

	"github.com/NavaneethWKT/CapStone_GO_Lang/protoc"
	"github.com/NavaneethWKT/CapStone_GO_Lang/server/internal/errors"
	"github.com/NavaneethWKT/CapStone_GO_Lang/server/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// BuyVoucher handles voucher purchase
func (h *PaymentHandler) BuyVoucher(ctx context.Context, req *protoc.BuyVoucherRequest) (*protoc.BuyVoucherResponse, error) {
	userID := int(req.GetUserId())
	voucherID := int(req.GetVoucherId())

	// Call service
	transaction, err := h.paymentService.BuyVoucher(userID, voucherID)
	if err != nil {
		return nil, h.handleError(err)
	}

	// Convert domain model to gRPC message
	pbTransaction := &protoc.Transaction{
		Id:              int32(transaction.ID),
		UserId:          int32(transaction.UserID),
		Amount:          transaction.Amount,
		TransactionType: string(transaction.TransactionType),
		PaymentStatus:   string(transaction.PaymentStatus),
		CreatedAt:       transaction.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       transaction.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if transaction.VoucherID != nil {
		vID := int32(*transaction.VoucherID)
		pbTransaction.VoucherId = vID
	}

	if transaction.PaymentTxnID != nil {
		pbTransaction.PaymentTxnId = *transaction.PaymentTxnID
	}

	return &protoc.BuyVoucherResponse{
		Transaction: pbTransaction,
		Message:     "Voucher purchased successfully",
	}, nil
}

// handleError converts application errors to gRPC status errors
func (h *PaymentHandler) handleError(err error) error {
	switch err {
	case errors.ErrUserNotFound:
		return status.Error(codes.NotFound, err.Error())
	case errors.ErrVoucherNotFound:
		return status.Error(codes.NotFound, err.Error())
	case errors.ErrVoucherOutOfStock:
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.ErrVoucherExpired:
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.ErrInsufficientBalance:
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.ErrPaymentFailed:
		return status.Error(codes.Internal, err.Error())
	case errors.ErrInvalidUserID, errors.ErrInvalidVoucherID:
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}

