package user

import "backend_path/internal/domain"

type Repository interface {
	Create(user *domain.User) error
	GetByID(id int) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id int) error
	GetAll() ([]*domain.User, error)
}
