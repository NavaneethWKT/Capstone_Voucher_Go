package handler

import (
	"context"
	"strings"

	"github.com/NavaneethWKT/CapStone_GO_Lang/protoc"
	"github.com/NavaneethWKT/CapStone_GO_Lang/server/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WalletHandler struct {
	walletService *service.WalletService
}

// NewWalletHandler creates a new wallet handler
func NewWalletHandler(walletService *service.WalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

// GetBalance retrieves wallet balance for a user
func (h *WalletHandler) GetBalance(ctx context.Context, req *protoc.GetBalanceRequest) (*protoc.GetBalanceResponse, error) {
	userID := int(req.GetUserId())

	// Call service
	balance, err := h.walletService.GetBalance(userID)
	if err != nil {
		return nil, h.handleError(err)
	}

	return &protoc.GetBalanceResponse{
		Balance: balance,
	}, nil
}

// handleError converts application errors to gRPC status errors
func (h *WalletHandler) handleError(err error) error {
	errMsg := err.Error()
	switch {
	case strings.Contains(errMsg, "user not found"):
		return status.Error(codes.NotFound, errMsg)
	case strings.Contains(errMsg, "invalid user ID"):
		return status.Error(codes.InvalidArgument, errMsg)
	case strings.Contains(errMsg, "wallet not found"):
		return status.Error(codes.NotFound, errMsg)
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}

