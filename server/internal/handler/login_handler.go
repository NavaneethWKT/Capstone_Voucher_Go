package handler

import (
	"context"
	"strings"
	"time"

	"github.com/NavaneethWKT/CapStone_GO_Lang/protoc"
	"github.com/NavaneethWKT/CapStone_GO_Lang/server/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LoginHandler struct {
	userService *service.UserService
}

// NewLoginHandler creates a new login handler
func NewLoginHandler(userService *service.UserService) *LoginHandler {
	return &LoginHandler{
		userService: userService,
	}
}

// Login authenticates a user with email and password
func (h *LoginHandler) Login(ctx context.Context, req *protoc.LoginRequest) (*protoc.LoginResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()

	if email == "" || password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	// Call service
	user, err := h.userService.Login(email, password)
	if err != nil {
		return nil, h.handleError(err)
	}

	// Convert domain model to gRPC message
	pbUser := &protoc.User{
		Id:        int32(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}

	return &protoc.LoginResponse{
		User:    pbUser,
		Message: "Login successful",
	}, nil
}

// handleError converts application errors to gRPC status errors
func (h *LoginHandler) handleError(err error) error {
	errMsg := err.Error()
	switch {
	case strings.Contains(errMsg, "invalid email or password"):
		return status.Error(codes.Unauthenticated, errMsg)
	case strings.Contains(errMsg, "user not found"):
		return status.Error(codes.NotFound, errMsg)
	case strings.Contains(errMsg, "invalid user ID"):
		return status.Error(codes.InvalidArgument, errMsg)
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}

