package handlers

import (
	"net/http"
)

func Health(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte("Healthy"))
}
