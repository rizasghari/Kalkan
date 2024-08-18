package handlers

import "net/http"

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Helath(w http.ResponseWriter, r *http.Request) {	
	w.WriteHeader(http.StatusOK)
}