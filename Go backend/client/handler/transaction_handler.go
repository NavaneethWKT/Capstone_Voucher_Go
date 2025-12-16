package handler

import (
	"net/http"
	"strconv"

	"github.com/NavaneethWKT/CapStone_GO_Lang/client/service"
	"github.com/NavaneethWKT/CapStone_GO_Lang/protoc"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TransactionHandler struct {
	grpcClient *service.GRPCClient
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(grpcClient *service.GRPCClient) *TransactionHandler {
	return &TransactionHandler{
		grpcClient: grpcClient,
	}
}

// ListTransactions handles GET /api/v1/transactions/:user_id
func (h *TransactionHandler) ListTransactions(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid user_id",
		})
		return
	}

	// Build gRPC request
	req := &protoc.ListTransactionsRequest{
		UserId: int32(userID),
	}

	// Call gRPC server
	client := h.grpcClient.GetVoucherClient()
	resp, err := client.ListTransactions(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"transactions": resp.GetTransactions(),
	})
}

// handleError converts gRPC errors to HTTP responses
func (h *TransactionHandler) handleError(c *gin.Context, err error) {
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
	default:
		httpStatus = http.StatusInternalServerError
	}

	c.JSON(httpStatus, gin.H{
		"success": false,
		"error":   st.Message(),
	})
}

