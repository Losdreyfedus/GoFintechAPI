package balance

import (
	"errors"
	"time"

	"backend_path/internal/domain"
	"backend_path/pkg/logger"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) BalanceService {
	return &service{repo: repo}
}

func (s *service) UpdateBalance(userID int, amount float64) error {
	// Get current balance
	balance, err := s.repo.GetByUserID(userID)
	if err != nil {
		// Create new balance if not exists
		balance = &domain.Balance{
			UserID:        userID,
			Amount:        0,
			LastUpdatedAt: time.Now(),
		}
	}

	// Thread-safe balance update
	balance.Update(amount)

	// Save to database
	if err := s.repo.Update(balance); err != nil {
		logger.Error("Failed to update balance", err, map[string]interface{}{
			"user_id": userID,
			"amount":  amount,
		})
		return err
	}

	logger.Info("Balance updated successfully", map[string]interface{}{
		"user_id": userID,
		"amount":  amount,
	})

	return nil
}

func (s *service) GetCurrentBalance(userID int) (float64, error) {
	balance, err := s.repo.GetByUserID(userID)
	if err != nil {
		logger.Error("Failed to get current balance", err, map[string]interface{}{
			"user_id": userID,
		})
		return 0, err
	}

	return balance.GetAmount(), nil
}

func (s *service) GetHistoricalBalance(userID int, atTime string) (float64, error) {
	// Parse the time string
	_, err := time.Parse(time.RFC3339, atTime)
	if err != nil {
		return 0, errors.New("invalid time format, expected RFC3339")
	}

	// For now, return current balance as historical balance
	// In a real implementation, you would query historical balance data
	balance, err := s.repo.GetByUserID(userID)
	if err != nil {
		logger.Error("Failed to get historical balance", err, map[string]interface{}{
			"user_id": userID,
			"at_time": atTime,
		})
		return 0, err
	}

	logger.Info("Historical balance retrieved", map[string]interface{}{
		"user_id": userID,
		"at_time": atTime,
		"balance": balance.GetAmount(),
	})

	return balance.GetAmount(), nil
}
