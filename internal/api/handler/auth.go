package handler

import (
	"encoding/json"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call userService.Register, return response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse request, call userService.Authenticate, return token
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": "dummy-token"})
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement token refresh logic
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": "refreshed-token"})
}
