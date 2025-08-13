package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend_path/internal/api/dto"
	"backend_path/internal/user"
	"backend_path/pkg/logger"

	"github.com/go-chi/chi/v5"
)

var userService user.UserService

// SetUserService sets the user service dependency
func SetUserService(service user.UserService) {
	userService = service
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := userService.GetAllUsers()
	if err != nil {
		logger.Error("Failed to get users", err, nil)
		respondWithError(w, http.StatusInternalServerError, "Failed to get users", err)
		return
	}

	response := make([]dto.UserResponse, len(users))
	for i, u := range users {
		response[i] = dto.UserResponse{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			Role:      u.Role,
			CreatedAt: u.CreatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	user, err := userService.GetByID(userID)
	if err != nil {
		logger.Error("Failed to get user", err, map[string]interface{}{
			"user_id": userID,
		})
		respondWithError(w, http.StatusNotFound, "User not found", err)
		return
	}

	response := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Update user
	user, err := userService.UpdateUser(userID, req.Username, req.Email, req.Role)
	if err != nil {
		logger.Error("Failed to update user", err, map[string]interface{}{
			"user_id": userID,
		})
		respondWithError(w, http.StatusInternalServerError, "Failed to update user", err)
		return
	}

	response := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	if err := userService.DeleteUser(userID); err != nil {
		logger.Error("Failed to delete user", err, map[string]interface{}{
			"user_id": userID,
		})
		respondWithError(w, http.StatusInternalServerError, "Failed to delete user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
