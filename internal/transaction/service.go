package transaction

import (
	"errors"
	"time"

	"backend_path/internal/domain"
	"backend_path/pkg/logger"
)

type service struct {
	repo           Repository
	balanceService interface {
		UpdateBalance(userID int, amount float64) error
	}
}

func NewService(repo Repository, balanceService interface {
	UpdateBalance(userID int, amount float64) error
}) TransactionService {
	return &service{repo: repo, balanceService: balanceService}
}

func (s *service) ProcessCredit(userID int, amount float64) (*domain.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("credit amount must be positive")
	}

	// Create credit transaction
	tx := &domain.Transaction{
		FromUserID: -1, // System credit (will be set to NULL in database)
		ToUserID:   userID,
		Amount:     amount,
		Type:       "credit",
		Status:     domain.StatusPending,
		CreatedAt:  time.Now(),
	}

	// Save transaction
	if err := s.repo.Create(tx); err != nil {
		logger.Error("Failed to create credit transaction", err, map[string]interface{}{
			"user_id": userID,
			"amount":  amount,
		})
		return nil, err
	}

	// Update status to completed
	tx.SetStatus(domain.StatusCompleted)
	if err := s.repo.Update(tx); err != nil {
		logger.Error("Failed to update credit transaction status", err, map[string]interface{}{
			"transaction_id": tx.ID,
		})
		return nil, err
	}

	// Update user balance
	if err := s.balanceService.UpdateBalance(userID, amount); err != nil {
		logger.Error("Failed to update user balance after credit", err, map[string]interface{}{
			"user_id": userID,
			"amount":  amount,
		})
		// Don't fail the transaction if balance update fails
	}

	logger.Info("Credit transaction processed successfully", map[string]interface{}{
		"transaction_id": tx.ID,
		"user_id":        userID,
		"amount":         amount,
	})

	return tx, nil
}

func (s *service) ProcessDebit(userID int, amount float64) (*domain.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("debit amount must be positive")
	}

	// Create debit transaction
	tx := &domain.Transaction{
		FromUserID: userID,
		ToUserID:   -1, // System debit (will be set to NULL in database)
		Amount:     amount,
		Type:       "debit",
		Status:     domain.StatusPending,
		CreatedAt:  time.Now(),
	}

	// Save transaction
	if err := s.repo.Create(tx); err != nil {
		logger.Error("Failed to create debit transaction", err, map[string]interface{}{
			"user_id": userID,
			"amount":  amount,
		})
		return nil, err
	}

	// Update status to completed
	tx.SetStatus(domain.StatusCompleted)
	if err := s.repo.Update(tx); err != nil {
		logger.Error("Failed to update debit transaction status", err, map[string]interface{}{
			"transaction_id": tx.ID,
		})
		return nil, err
	}

	// Update user balance (subtract amount)
	if err := s.balanceService.UpdateBalance(userID, -amount); err != nil {
		logger.Error("Failed to update user balance after debit", err, map[string]interface{}{
			"user_id": userID,
			"amount":  amount,
		})
		// Don't fail the transaction if balance update fails
	}

	logger.Info("Debit transaction processed successfully", map[string]interface{}{
		"transaction_id": tx.ID,
		"user_id":        userID,
		"amount":         amount,
	})

	return tx, nil
}

func (s *service) ProcessTransfer(fromUserID, toUserID int, amount float64) (*domain.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("transfer amount must be positive")
	}

	if fromUserID == toUserID {
		return nil, errors.New("cannot transfer to same user")
	}

	// Create transfer transaction
	tx := &domain.Transaction{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Amount:     amount,
		Type:       "transfer",
		Status:     domain.StatusPending,
		CreatedAt:  time.Now(),
	}

	// Save transaction
	if err := s.repo.Create(tx); err != nil {
		logger.Error("Failed to create transfer transaction", err, map[string]interface{}{
			"from_user_id": fromUserID,
			"to_user_id":   toUserID,
			"amount":       amount,
		})
		return nil, err
	}

	// Update status to completed
	tx.SetStatus(domain.StatusCompleted)
	if err := s.repo.Update(tx); err != nil {
		logger.Error("Failed to update transfer transaction status", err, map[string]interface{}{
			"transaction_id": tx.ID,
		})
		return nil, err
	}

	// Update sender balance (subtract amount)
	if err := s.balanceService.UpdateBalance(fromUserID, -amount); err != nil {
		logger.Error("Failed to update sender balance after transfer", err, map[string]interface{}{
			"from_user_id": fromUserID,
			"amount":       amount,
		})
		// Don't fail the transaction if balance update fails
	}

	// Update receiver balance (add amount)
	if err := s.balanceService.UpdateBalance(toUserID, amount); err != nil {
		logger.Error("Failed to update receiver balance after transfer", err, map[string]interface{}{
			"to_user_id": toUserID,
			"amount":     amount,
		})
		// Don't fail the transaction if balance update fails
	}

	logger.Info("Transfer transaction processed successfully", map[string]interface{}{
		"transaction_id": tx.ID,
		"from_user_id":   fromUserID,
		"to_user_id":     toUserID,
		"amount":         amount,
	})

	return tx, nil
}

func (s *service) GetTransaction(id int) (*domain.Transaction, error) {
	tx, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error("Failed to get transaction", err, map[string]interface{}{
			"transaction_id": id,
		})
		return nil, err
	}

	return tx, nil
}

func (s *service) GetTransactionHistory(userID int) ([]*domain.Transaction, error) {
	transactions, err := s.repo.GetByUser(userID)
	if err != nil {
		logger.Error("Failed to get transaction history", err, map[string]interface{}{
			"user_id": userID,
		})
		return nil, err
	}

	return transactions, nil
}
