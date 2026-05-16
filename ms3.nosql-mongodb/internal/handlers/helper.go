package handlers

import (
	"encoding/json"
	"net/http"
)

func respondJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func respondError(w http.ResponseWriter, code int, msg string) {
	respondJSON(w, code, map[string]string{"error": msg})
}
