package handler

import (
	"context"
	"strings"

	"github.com/NavaneethWKT/CapStone_GO_Lang/protoc"
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
	errMsg := err.Error()
	switch {
	case strings.Contains(errMsg, "user not found"):
		return status.Error(codes.NotFound, errMsg)
	case strings.Contains(errMsg, "voucher not found"):
		return status.Error(codes.NotFound, errMsg)
	case strings.Contains(errMsg, "voucher out of stock"):
		return status.Error(codes.FailedPrecondition, errMsg)
	case strings.Contains(errMsg, "voucher expired"):
		return status.Error(codes.FailedPrecondition, errMsg)
	case strings.Contains(errMsg, "insufficient wallet balance"):
		return status.Error(codes.FailedPrecondition, errMsg)
	case strings.Contains(errMsg, "payment processing failed"):
		return status.Error(codes.Internal, errMsg)
	case strings.Contains(errMsg, "invalid user ID") || strings.Contains(errMsg, "invalid voucher ID"):
		return status.Error(codes.InvalidArgument, errMsg)
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}

