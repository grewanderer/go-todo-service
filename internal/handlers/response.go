package handlers

import (
	"encoding/json"
	"net/http"
)

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, r *http.Request, status int, message string) {
	response := map[string]string{"error": message}
	if r != nil {
		if requestID, ok := RequestIDFromContext(r.Context()); ok {
			response["request_id"] = requestID
		}
	}
	respondJSON(w, status, response)
}
