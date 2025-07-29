package handler

import (
	"encoding/json"
	"net/http"
)

func CurrentBalance(w http.ResponseWriter, r *http.Request) {
	// TODO: Return current balance
	json.NewEncoder(w).Encode(map[string]interface{}{"balance": 0})
}

func HistoricalBalance(w http.ResponseWriter, r *http.Request) {
	// TODO: Return historical balance
	json.NewEncoder(w).Encode([]interface{}{})
}

func BalanceAtTime(w http.ResponseWriter, r *http.Request) {
	// TODO: Return balance at a specific time
	json.NewEncoder(w).Encode(map[string]interface{}{"balance": 0})
}
