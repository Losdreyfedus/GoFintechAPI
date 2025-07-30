package transaction

import "backend_path/internal/domain"

// TransactionService provides transaction-related operations
type TransactionService interface {
	ProcessCredit(userID int, amount float64) (*domain.Transaction, error)
	ProcessDebit(userID int, amount float64) (*domain.Transaction, error)
	ProcessTransfer(fromUserID, toUserID int, amount float64) (*domain.Transaction, error)
	GetTransaction(id int) (*domain.Transaction, error)
	GetTransactionHistory(userID int) ([]*domain.Transaction, error)
}
