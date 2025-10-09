package handlers

import (
	"net/http"
)

type Handler struct{}

func NewHandler() Handler {
	return Handler{}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Implement your request handling logic here
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}
