package server

import (
	"encoding/json"
	"net/http"
)

// utility function to write content out to JSON
func writeJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(v)
}
