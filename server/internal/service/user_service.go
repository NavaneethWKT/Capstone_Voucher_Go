package service

import (
	"github.com/NavaneethWKT/CapStone_GO_Lang/server/internal/errors"
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
		return errors.ErrInvalidUserID
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return errors.WrapError(err, "failed to get user")
	}

	if user == nil {
		return errors.ErrUserNotFound
	}

	return nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(userID int) (*model.User, error) {
	if userID <= 0 {
		return nil, errors.ErrInvalidUserID
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.WrapError(err, "failed to get user")
	}

	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

