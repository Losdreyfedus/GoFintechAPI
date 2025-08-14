package domain

import (
	"encoding/json"
	"errors"
	"regexp"
	"sync"
	"time"
)

// User represents a user in the system
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) Validate() error {
	if u.Username == "" {
		return errors.New("username is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	if !regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`).MatchString(u.Email) {
		return errors.New("invalid email format")
	}
	if u.Role == "" {
		return errors.New("role is required")
	}
	return nil
}

func (u *User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(u),
	})
}

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	StatusPending    TransactionStatus = "pending"
	StatusCompleted  TransactionStatus = "completed"
	StatusFailed     TransactionStatus = "failed"
	StatusRolledBack TransactionStatus = "rolled_back"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID         int               `json:"id"`
	FromUserID int               `json:"from_user_id"`
	ToUserID   int               `json:"to_user_id"`
	Amount     float64           `json:"amount"`
	Type       string            `json:"type"`
	Status     TransactionStatus `json:"status"`
	CreatedAt  time.Time         `json:"created_at"`
}

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

// Balance represents a user's balance
type Balance struct {
	UserID        int       `json:"user_id"`
	Amount        float64   `json:"amount"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
	mu            sync.RWMutex
}

func (b *Balance) Update(amount float64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Amount += amount
	b.LastUpdatedAt = time.Now()
}

func (b *Balance) GetAmount() float64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Amount
}

func (b *Balance) MarshalJSON() ([]byte, error) {
	type Alias Balance
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(b),
	})
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID         int       `json:"id"`
	EntityType string    `json:"entity_type"`
	EntityID   int       `json:"entity_id"`
	Action     string    `json:"action"`
	Details    string    `json:"details"`
	CreatedAt  time.Time `json:"created_at"`
}
