package balance

import "backend_path/internal/domain"

type Repository interface {
	GetByUserID(userID int) (*domain.Balance, error)
	Update(balance *domain.Balance) error
}
