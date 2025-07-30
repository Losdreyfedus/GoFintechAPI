package user

import "backend_path/internal/domain"

// UserService provides user-related operations
type UserService interface {
	Register(user *domain.User, password string) error
	Authenticate(email, password string) (*domain.User, error)
	Authorize(user *domain.User, role string) bool
	GetByID(id int) (*domain.User, error)
}
