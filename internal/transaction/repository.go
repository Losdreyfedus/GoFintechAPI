package transaction

import "backend_path/internal/domain"

type Repository interface {
	Create(tx *domain.Transaction) error
	GetByID(id int) (*domain.Transaction, error)
	GetByUser(userID int) ([]*domain.Transaction, error)
	Update(tx *domain.Transaction) error
}
