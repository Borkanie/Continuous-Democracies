package handlers

import (
	"encoding/json"
	"net/http"
)

func writeJSON(responseWriter http.ResponseWriter, status int, payload any) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(status)
	json.NewEncoder(responseWriter).Encode(payload)
}

func writeNotFound(responseWriter http.ResponseWriter, msg string) {
	http.Error(responseWriter, msg, http.StatusNotFound)
}

func writeBadRequest(responseWriter http.ResponseWriter, msg string) {
	http.Error(responseWriter, msg, http.StatusBadRequest)
}

func writeInternalError(responseWriter http.ResponseWriter, err error) {
	http.Error(responseWriter, "internal server error", http.StatusInternalServerError)
}
