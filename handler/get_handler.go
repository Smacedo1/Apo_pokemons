package handler

import (
	"encoding/json"
	"gos/service"
	"net/http"
)

type GetHandler struct {
	service service.Service
}

func NewGetHandler(s service.Service) *GetHandler {
	return &GetHandler{service: s}
}

func (h *handler) GetType(w http.ResponseWriter, r *http.Request) {
	types, err := h.service.GetType(r.Context())
	if err != nil {
		http.Error(w, "error obteniendo tipos: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(types)
	
}

