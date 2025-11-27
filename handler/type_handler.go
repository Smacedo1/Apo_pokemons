package handler

import ()

func (h *handler) GetTypes(w http.ResponseWriter, r *http.Request) { 
	types,  err := h.service.GetType(r.Context())
	if err != nil {
		http.Error(w, "error obteniendo tipos: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(types)
}

