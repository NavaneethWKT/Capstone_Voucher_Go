package handler

import (
	"context"

	"github.com/NavaneethWKT/CapStone_GO_Lang/protoc"
	"github.com/NavaneethWKT/CapStone_GO_Lang/server/internal/service"
)

// VoucherServiceHandler implements all gRPC service methods
type VoucherServiceHandler struct {
	protoc.UnimplementedVoucherServiceServer
	voucherHandler     *VoucherHandler
	paymentHandler     *PaymentHandler
	walletHandler      *WalletHandler
	transactionHandler *TransactionHandler
}

// NewVoucherServiceHandler creates a new combined handler
func NewVoucherServiceHandler(
	voucherService *service.VoucherService,
	paymentService *service.PaymentService,
	walletService *service.WalletService,
	transactionService *service.TransactionService,
) *VoucherServiceHandler {
	return &VoucherServiceHandler{
		voucherHandler:     NewVoucherHandler(voucherService),
		paymentHandler:     NewPaymentHandler(paymentService),
		walletHandler:      NewWalletHandler(walletService),
		transactionHandler: NewTransactionHandler(transactionService),
	}
}

// Search delegates to VoucherHandler
func (h *VoucherServiceHandler) Search(ctx context.Context, req *protoc.SearchRequest) (*protoc.SearchResponse, error) {
	return h.voucherHandler.Search(ctx, req)
}

// BuyVoucher delegates to PaymentHandler
func (h *VoucherServiceHandler) BuyVoucher(ctx context.Context, req *protoc.BuyVoucherRequest) (*protoc.BuyVoucherResponse, error) {
	return h.paymentHandler.BuyVoucher(ctx, req)
}

// GetBalance delegates to WalletHandler
func (h *VoucherServiceHandler) GetBalance(ctx context.Context, req *protoc.GetBalanceRequest) (*protoc.GetBalanceResponse, error) {
	return h.walletHandler.GetBalance(ctx, req)
}

// ListTransactions delegates to TransactionHandler
func (h *VoucherServiceHandler) ListTransactions(ctx context.Context, req *protoc.ListTransactionsRequest) (*protoc.ListTransactionsResponse, error) {
	return h.transactionHandler.ListTransactions(ctx, req)
}

