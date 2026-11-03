package handler

import (
	"encoding/json"
	"net/http"
)

func WriteError(w http.ResponseWriter, status int, data any) {
	WriteResponse(w, status, data)
}

func WriteResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
