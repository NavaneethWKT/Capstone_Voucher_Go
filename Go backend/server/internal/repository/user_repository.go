package repository

import (
	"database/sql"

	"github.com/NavaneethWKT/CapStone_GO_Lang/server/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(id int) (*model.User, error) {
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	
	user := &model.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}
	
	return user, nil
}

// Login checks if email and password match
func (r *UserRepository) Login(email, password string) (*model.User, error) {
	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE email = $1 AND password = $2
	`
	
	user := &model.User{}
	err := r.db.QueryRow(query, email, password).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Invalid credentials
		}
		return nil, err
	}
	
	return user, nil
}

