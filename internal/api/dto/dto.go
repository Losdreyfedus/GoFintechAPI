package dto

import "time"

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,oneof=user admin"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshRequest represents token refresh request
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int       `json:"expires_in"`
	User         UserInfo  `json:"user"`
	Timestamp    time.Time `json:"timestamp"`
}

// UserInfo represents user information in responses
type UserInfo struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// CreditRequest represents credit transaction request
type CreditRequest struct {
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Description string  `json:"description,omitempty"`
}

// DebitRequest represents debit transaction request
type DebitRequest struct {
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Description string  `json:"description,omitempty"`
}

// TransactionRequest represents transaction request
type TransactionRequest struct {
	UserID      int     `json:"user_id" validate:"required"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Description string  `json:"description" validate:"required"`
	Currency    string  `json:"currency" validate:"required,len=3"`
	Reference   string  `json:"reference,omitempty"`
	Category    string  `json:"category,omitempty"`
}

// TransferRequest represents transfer request
type TransferRequest struct {
	ToUserID    int     `json:"to_user_id" validate:"required"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Description string  `json:"description,omitempty"`
}

// TransactionResponse represents transaction response
type TransactionResponse struct {
	ID         int       `json:"id"`
	FromUserID int       `json:"from_user_id"`
	ToUserID   int       `json:"to_user_id"`
	Amount     float64   `json:"amount"`
	Type       string    `json:"type"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

// TransferResponse represents transfer response
type TransferResponse struct {
	TransactionID   string    `json:"transaction_id"`
	Status          string    `json:"status"`
	Amount          float64   `json:"amount"`
	FromUserBalance float64   `json:"from_user_balance"`
	ToUserBalance   float64   `json:"to_user_balance"`
	Description     string    `json:"description"`
	Currency        string    `json:"currency"`
	Timestamp       time.Time `json:"timestamp"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// BalanceResponse represents balance response
type BalanceResponse struct {
	UserID  int     `json:"user_id"`
	Amount  float64 `json:"amount"`
	Type    string  `json:"type"`
	Updated string  `json:"updated"`
}

// SuccessResponse represents success response
type SuccessResponse struct {
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}
