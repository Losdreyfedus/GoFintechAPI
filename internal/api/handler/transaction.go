package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend_path/internal/api/dto"
	"backend_path/internal/transaction"
	"backend_path/pkg/logger"

	"github.com/go-chi/chi/v5"
)

var transactionService transaction.TransactionService

// SetTransactionService sets the transaction service dependency
func SetTransactionService(service transaction.TransactionService) {
	transactionService = service
}

func Credit(w http.ResponseWriter, r *http.Request) {
	var req dto.CreditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if req.Amount <= 0 {
		respondWithError(w, http.StatusBadRequest, "Amount must be positive", nil)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := getUserIDFromContext(r)
	if userID == 0 {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	// Process credit transaction
	tx, err := transactionService.ProcessCredit(userID, req.Amount)
	if err != nil {
		logger.Error("Failed to process credit", err, map[string]interface{}{
			"user_id": userID,
			"amount":  req.Amount,
		})
		respondWithError(w, http.StatusInternalServerError, "Failed to process credit", err)
		return
	}

	response := dto.TransactionResponse{
		ID:         tx.ID,
		FromUserID: tx.FromUserID,
		ToUserID:   tx.ToUserID,
		Amount:     tx.Amount,
		Type:       tx.Type,
		Status:     string(tx.Status),
		CreatedAt:  tx.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func Debit(w http.ResponseWriter, r *http.Request) {
	var req dto.DebitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if req.Amount <= 0 {
		respondWithError(w, http.StatusBadRequest, "Amount must be positive", nil)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := getUserIDFromContext(r)
	if userID == 0 {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	// Process debit transaction
	tx, err := transactionService.ProcessDebit(userID, req.Amount)
	if err != nil {
		logger.Error("Failed to process debit", err, map[string]interface{}{
			"user_id": userID,
			"amount":  req.Amount,
		})
		respondWithError(w, http.StatusInternalServerError, "Failed to process debit", err)
		return
	}

	response := dto.TransactionResponse{
		ID:         tx.ID,
		FromUserID: tx.FromUserID,
		ToUserID:   tx.ToUserID,
		Amount:     tx.Amount,
		Type:       tx.Type,
		Status:     string(tx.Status),
		CreatedAt:  tx.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func Transfer(w http.ResponseWriter, r *http.Request) {
	var req dto.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if req.Amount <= 0 {
		respondWithError(w, http.StatusBadRequest, "Amount must be positive", nil)
		return
	}

	if req.ToUserID <= 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid recipient user ID", nil)
		return
	}

	// Get user ID from context (set by auth middleware)
	fromUserID := getUserIDFromContext(r)
	if fromUserID == 0 {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	// Process transfer transaction
	tx, err := transactionService.ProcessTransfer(fromUserID, req.ToUserID, req.Amount)
	if err != nil {
		logger.Error("Failed to process transfer", err, map[string]interface{}{
			"from_user_id": fromUserID,
			"to_user_id":   req.ToUserID,
			"amount":       req.Amount,
		})
		respondWithError(w, http.StatusInternalServerError, "Failed to process transfer", err)
		return
	}

	response := dto.TransactionResponse{
		ID:         tx.ID,
		FromUserID: tx.FromUserID,
		ToUserID:   tx.ToUserID,
		Amount:     tx.Amount,
		Type:       tx.Type,
		Status:     string(tx.Status),
		CreatedAt:  tx.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func TransactionHistory(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := getUserIDFromContext(r)
	if userID == 0 {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	// Get transaction history
	transactions, err := transactionService.GetTransactionHistory(userID)
	if err != nil {
		logger.Error("Failed to get transaction history", err, map[string]interface{}{
			"user_id": userID,
		})
		respondWithError(w, http.StatusInternalServerError, "Failed to get transaction history", err)
		return
	}

	// Convert to response format
	response := make([]dto.TransactionResponse, len(transactions))
	for i, tx := range transactions {
		response[i] = dto.TransactionResponse{
			ID:         tx.ID,
			FromUserID: tx.FromUserID,
			ToUserID:   tx.ToUserID,
			Amount:     tx.Amount,
			Type:       tx.Type,
			Status:     string(tx.Status),
			CreatedAt:  tx.CreatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetTransaction(w http.ResponseWriter, r *http.Request) {
	// Get transaction ID from URL
	transactionIDStr := chi.URLParam(r, "id")
	transactionID, err := strconv.Atoi(transactionIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid transaction ID", err)
		return
	}

	// Get transaction
	tx, err := transactionService.GetTransaction(transactionID)
	if err != nil {
		logger.Error("Failed to get transaction", err, map[string]interface{}{
			"transaction_id": transactionID,
		})
		respondWithError(w, http.StatusNotFound, "Transaction not found", err)
		return
	}

	response := dto.TransactionResponse{
		ID:         tx.ID,
		FromUserID: tx.FromUserID,
		ToUserID:   tx.ToUserID,
		Amount:     tx.Amount,
		Type:       tx.Type,
		Status:     string(tx.Status),
		CreatedAt:  tx.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Helper function to get user ID from context
func getUserIDFromContext(r *http.Request) int {
	userID := r.Context().Value("user")
	if userID == nil {
		return 0
	}

	if id, ok := userID.(int); ok {
		return id
	}

	return 0
}
