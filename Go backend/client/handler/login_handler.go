package handler

import (
	"net/http"

	"github.com/NavaneethWKT/CapStone_GO_Lang/client/service"
	"github.com/NavaneethWKT/CapStone_GO_Lang/protoc"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LoginHandler struct {
	grpcClient *service.GRPCClient
}

// NewLoginHandler creates a new login handler
func NewLoginHandler(grpcClient *service.GRPCClient) *LoginHandler {
	return &LoginHandler{
		grpcClient: grpcClient,
	}
}

// Login handles POST /api/v1/auth/login
func (h *LoginHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid request body",
		})
		return
	}

	// Build gRPC request
	grpcReq := &protoc.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	// Call gRPC server
	client := h.grpcClient.GetVoucherClient()
	resp, err := client.Login(c.Request.Context(), grpcReq)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    resp.GetUser(),
		"message": resp.GetMessage(),
	})
}

// handleError converts gRPC errors to HTTP responses
func (h *LoginHandler) handleError(c *gin.Context, err error) {
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
	case codes.Unauthenticated:
		httpStatus = http.StatusUnauthorized
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

