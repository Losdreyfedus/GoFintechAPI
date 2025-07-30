package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"backend_path/internal/api/dto"
	"backend_path/internal/domain"
	"backend_path/internal/user"
	"backend_path/pkg/jwt"
	"backend_path/pkg/logger"
)

type AuthHandler struct {
	userService user.UserService
	jwtService  *jwt.JWTService
}

func NewAuthHandler(userService user.UserService, jwtService *jwt.JWTService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Create user domain object
	user := &domain.User{
		Username: req.Username,
		Email:    req.Email,
		Role:     req.Role,
	}

	// Register user
	if err := h.userService.Register(user, req.Password); err != nil {
		logger.Error("Failed to register user", err, map[string]interface{}{
			"email": req.Email,
		})
		respondWithError(w, http.StatusInternalServerError, "Failed to register user", err)
		return
	}

	// Generate tokens
	token, err := h.jwtService.GenerateToken(user.ID, user.Username, user.Email, user.Role)
	if err != nil {
		logger.Error("Failed to generate token", err, nil)
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token", err)
		return
	}

	refreshToken, err := h.jwtService.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		logger.Error("Failed to generate refresh token", err, nil)
		respondWithError(w, http.StatusInternalServerError, "Failed to generate refresh token", err)
		return
	}

	response := dto.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresIn:    3600, // 1 hour
		User: dto.UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Authenticate user
	user, err := h.userService.Authenticate(req.Email, req.Password)
	if err != nil {
		logger.Error("Authentication failed", err, map[string]interface{}{
			"email": req.Email,
		})
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials", err)
		return
	}

	// Generate tokens
	token, err := h.jwtService.GenerateToken(user.ID, user.Username, user.Email, user.Role)
	if err != nil {
		logger.Error("Failed to generate token", err, nil)
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token", err)
		return
	}

	refreshToken, err := h.jwtService.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		logger.Error("Failed to generate refresh token", err, nil)
		respondWithError(w, http.StatusInternalServerError, "Failed to generate refresh token", err)
		return
	}

	response := dto.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresIn:    3600, // 1 hour
		User: dto.UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate refresh token
	claims, err := h.jwtService.ValidateToken(req.RefreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	// Get user
	user, err := h.userService.GetByID(claims.UserID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User not found", err)
		return
	}

	// Generate new tokens
	token, err := h.jwtService.GenerateToken(user.ID, user.Username, user.Email, user.Role)
	if err != nil {
		logger.Error("Failed to generate token", err, nil)
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token", err)
		return
	}

	refreshToken, err := h.jwtService.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		logger.Error("Failed to generate refresh token", err, nil)
		respondWithError(w, http.StatusInternalServerError, "Failed to generate refresh token", err)
		return
	}

	response := dto.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresIn:    3600, // 1 hour
		User: dto.UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Helper function to respond with error
func respondWithError(w http.ResponseWriter, statusCode int, message string, err error) {
	response := dto.ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	}

	if err != nil {
		response.Details = map[string]string{
			"error": err.Error(),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
