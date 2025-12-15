package handler

import (
	"fmt"
	"net/http"

	"github.com/NavaneethWKT/CapStone_GO_Lang/client/service"
	"github.com/NavaneethWKT/CapStone_GO_Lang/protoc"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type VoucherHandler struct {
	grpcClient *service.GRPCClient
}

// NewVoucherHandler creates a new voucher handler
func NewVoucherHandler(grpcClient *service.GRPCClient) *VoucherHandler {
	return &VoucherHandler{
		grpcClient: grpcClient,
	}
}

// Search handles GET /api/v1/vouchers/search
func (h *VoucherHandler) Search(c *gin.Context) {
	// Extract query parameters
	category := c.Query("category")
	minPriceStr := c.Query("min_price")
	maxPriceStr := c.Query("max_price")

	// Build gRPC request
	req := &protoc.SearchRequest{
		Category: category,
	}

	// Parse min_price if provided
	if minPriceStr != "" {
		var minPrice float64
		if _, err := fmt.Sscanf(minPriceStr, "%f", &minPrice); err == nil {
			req.MinPrice = minPrice
		}
	}

	// Parse max_price if provided
	if maxPriceStr != "" {
		var maxPrice float64
		if _, err := fmt.Sscanf(maxPriceStr, "%f", &maxPrice); err == nil {
			req.MaxPrice = maxPrice
		}
	}

	// Call gRPC server
	client := h.grpcClient.GetVoucherClient()
	resp, err := client.Search(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"vouchers": resp.GetVouchers(),
	})
}

// handleError converts gRPC errors to HTTP responses
func (h *VoucherHandler) handleError(c *gin.Context, err error) {
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

