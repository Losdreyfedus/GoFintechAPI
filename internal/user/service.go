package user

import (
	"backend_path/internal/domain"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) UserService {
	return &service{repo: repo}
}

func (s *service) Register(user *domain.User, password string) error {
	if err := user.Validate(); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hash)
	return s.repo.Create(user)
}

func (s *service) Authenticate(email, password string) (*domain.User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

func (s *service) Authorize(user *domain.User, role string) bool {
	return user.Role == role
}

// GetByID retrieves a user by ID
func (s *service) GetByID(id int) (*domain.User, error) {
	return s.repo.GetByID(id)
}
