package handler

import (
	"encoding/json"
	"net/http"

	"backend_path/internal/api/dto"
	"backend_path/internal/balance"
	"backend_path/pkg/logger"
)

var balanceService balance.BalanceService

// SetBalanceService sets the balance service dependency
func SetBalanceService(service balance.BalanceService) {
	balanceService = service
}

func CurrentBalance(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := getUserIDFromContext(r)
	if userID == 0 {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	// Get current balance
	currentBalance, err := balanceService.GetCurrentBalance(userID)
	if err != nil {
		logger.Error("Failed to get current balance", err, map[string]interface{}{
			"user_id": userID,
		})
		respondWithError(w, http.StatusInternalServerError, "Failed to get current balance", err)
		return
	}

	response := dto.BalanceResponse{
		UserID:  userID,
		Amount:  currentBalance,
		Type:    "current",
		Updated: "now",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func HistoricalBalance(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := getUserIDFromContext(r)
	if userID == 0 {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	// Get query parameters for date range (for future implementation)
	_ = r.URL.Query().Get("from")
	_ = r.URL.Query().Get("to")

	// For now, return current balance as historical data
	// In a real implementation, you would query historical balance data
	currentBalance, err := balanceService.GetCurrentBalance(userID)
	if err != nil {
		logger.Error("Failed to get historical balance", err, map[string]interface{}{
			"user_id": userID,
		})
		respondWithError(w, http.StatusInternalServerError, "Failed to get historical balance", err)
		return
	}

	response := []dto.BalanceResponse{
		{
			UserID:  userID,
			Amount:  currentBalance,
			Type:    "historical",
			Updated: "now",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func BalanceAtTime(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := getUserIDFromContext(r)
	if userID == 0 {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	// Get timestamp from query parameter
	atTime := r.URL.Query().Get("at")
	if atTime == "" {
		respondWithError(w, http.StatusBadRequest, "Missing 'at' parameter", nil)
		return
	}

	// Get balance at specific time
	balanceAtTime, err := balanceService.GetHistoricalBalance(userID, atTime)
	if err != nil {
		logger.Error("Failed to get balance at time", err, map[string]interface{}{
			"user_id": userID,
			"at_time": atTime,
		})
		respondWithError(w, http.StatusInternalServerError, "Failed to get balance at time", err)
		return
	}

	response := dto.BalanceResponse{
		UserID:  userID,
		Amount:  balanceAtTime,
		Type:    "at_time",
		Updated: atTime,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
