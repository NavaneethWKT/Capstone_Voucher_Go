package handler

import (
	"net/http"

	"github.com/NavaneethWKT/CapStone_GO_Lang/client/service"
	"github.com/NavaneethWKT/CapStone_GO_Lang/protoc"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentHandler struct {
	grpcClient *service.GRPCClient
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(grpcClient *service.GRPCClient) *PaymentHandler {
	return &PaymentHandler{
		grpcClient: grpcClient,
	}
}

// BuyVoucher handles POST /api/v1/vouchers/buy
func (h *PaymentHandler) BuyVoucher(c *gin.Context) {
	var req struct {
		UserID    int `json:"user_id" binding:"required"`
		VoucherID int `json:"voucher_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid request body",
		})
		return
	}

	// Build gRPC request
	grpcReq := &protoc.BuyVoucherRequest{
		UserId:    int32(req.UserID),
		VoucherId: int32(req.VoucherID),
	}

	// Call gRPC server
	client := h.grpcClient.GetVoucherClient()
	resp, err := client.BuyVoucher(c.Request.Context(), grpcReq)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"transaction": resp.GetTransaction(),
		"message":    resp.GetMessage(),
	})
}

// handleError converts gRPC errors to HTTP responses
func (h *PaymentHandler) handleError(c *gin.Context, err error) {
	st, ok := status.FromError(err)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "internal server error",
		})
		return
	}

	var httpStatus int
	switch st.Code() {
	case codes.NotFound:
		httpStatus = http.StatusNotFound
	case codes.InvalidArgument:
		httpStatus = http.StatusBadRequest
	case codes.FailedPrecondition:
		httpStatus = http.StatusPreconditionFailed
	default:
		httpStatus = http.StatusInternalServerError
	}

	c.JSON(httpStatus, gin.H{
		"success": false,
		"error":   st.Message(),
	})
}

