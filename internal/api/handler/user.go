package handler

import (
	"encoding/json"
	"net/http"
)

func ListUsers(w http.ResponseWriter, r *http.Request) {
	// TODO: Fetch and return list of users
	json.NewEncoder(w).Encode([]interface{}{})
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Fetch and return user by ID
	json.NewEncoder(w).Encode(map[string]interface{}{})
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Update user by ID
	w.WriteHeader(http.StatusNoContent)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Delete user by ID
	w.WriteHeader(http.StatusNoContent)
}
