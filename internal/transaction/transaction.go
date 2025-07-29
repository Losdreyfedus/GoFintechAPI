package domain

import (
	"encoding/json"
	"time"
)

type TransactionStatus string

type Transaction struct {
	ID         int               `json:"id"`
	FromUserID int               `json:"from_user_id"`
	ToUserID   int               `json:"to_user_id"`
	Amount     float64           `json:"amount"`
	Type       string            `json:"type"`
	Status     TransactionStatus `json:"status"`
	CreatedAt  time.Time         `json:"created_at"`
}

const (
	StatusPending    TransactionStatus = "pending"
	StatusCompleted  TransactionStatus = "completed"
	StatusFailed     TransactionStatus = "failed"
	StatusRolledBack TransactionStatus = "rolled_back"
)

func (t *Transaction) SetStatus(status TransactionStatus) {
	t.Status = status
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	type Alias Transaction
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	})
}
