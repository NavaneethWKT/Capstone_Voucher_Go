package service

import (
	"errors"
	"fmt"

	"github.com/NavaneethWKT/CapStone_GO_Lang/server/internal/model"
	"github.com/NavaneethWKT/CapStone_GO_Lang/server/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// ValidateUserExists checks if a user exists by ID
func (s *UserService) ValidateUserExists(userID int) error {
	if userID <= 0 {
		return errors.New("invalid user ID")
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return errors.New("user not found")
	}

	return nil
}

// Login authenticates a user with email and password
func (s *UserService) Login(email, password string) (*model.User, error) {
	if email == "" || password == "" {
		return nil, errors.New("invalid email or password")
	}

	user, err := s.userRepo.Login(email, password)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

