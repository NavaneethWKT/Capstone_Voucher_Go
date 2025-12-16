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

type WalletHandler struct {
	grpcClient *service.GRPCClient
}

// NewWalletHandler creates a new wallet handler
func NewWalletHandler(grpcClient *service.GRPCClient) *WalletHandler {
	return &WalletHandler{
		grpcClient: grpcClient,
	}
}

// GetBalance handles GET /api/v1/wallet/balance/:user_id
func (h *WalletHandler) GetBalance(c *gin.Context) {
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
	req := &protoc.GetBalanceRequest{
		UserId: int32(userID),
	}

	// Call gRPC server
	client := h.grpcClient.GetVoucherClient()
	resp, err := client.GetBalance(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"balance": resp.GetBalance(),
	})
}

// handleError converts gRPC errors to HTTP responses
func (h *WalletHandler) handleError(c *gin.Context, err error) {
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

