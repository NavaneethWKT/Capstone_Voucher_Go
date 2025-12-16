package handler

import (
	"context"
	"strings"

	"github.com/NavaneethWKT/CapStone_GO_Lang/protoc"
	"github.com/NavaneethWKT/CapStone_GO_Lang/server/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(transactionService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// ListTransactions retrieves all transactions for a user
func (h *TransactionHandler) ListTransactions(ctx context.Context, req *protoc.ListTransactionsRequest) (*protoc.ListTransactionsResponse, error) {
	userID := int(req.GetUserId())

	// Call service
	transactions, err := h.transactionService.ListTransactions(userID)
	if err != nil {
		return nil, h.handleError(err)
	}

	// Convert domain models to gRPC messages
	pbTransactions := make([]*protoc.Transaction, 0, len(transactions))
	for _, t := range transactions {
		pbTxn := &protoc.Transaction{
			Id:              int32(t.ID),
			UserId:          int32(t.UserID),
			Amount:          t.Amount,
			TransactionType: string(t.TransactionType),
			PaymentStatus:   string(t.PaymentStatus),
			CreatedAt:       t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:       t.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		if t.VoucherID != nil {
			vID := int32(*t.VoucherID)
			pbTxn.VoucherId = vID
		}

		if t.PaymentTxnID != nil {
			pbTxn.PaymentTxnId = *t.PaymentTxnID
		}

		pbTransactions = append(pbTransactions, pbTxn)
	}

	return &protoc.ListTransactionsResponse{
		Transactions: pbTransactions,
	}, nil
}

// handleError converts application errors to gRPC status errors
func (h *TransactionHandler) handleError(err error) error {
	errMsg := err.Error()
	switch {
	case strings.Contains(errMsg, "user not found"):
		return status.Error(codes.NotFound, errMsg)
	case strings.Contains(errMsg, "invalid user ID"):
		return status.Error(codes.InvalidArgument, errMsg)
	case strings.Contains(errMsg, "transaction not found"):
		return status.Error(codes.NotFound, errMsg)
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}

