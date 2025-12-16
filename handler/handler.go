package handler

import (
	"encoding/json"
	"gos/domain"
	"gos/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type handler struct { // guarda las dependencias que necesita para responder requests
	service service.Service
}

func NewHandler(s service.Service) *handler {
	return &handler{service: s}
}

func (h *handler) GetPokemons(w http.ResponseWriter, r *http.Request) {
	pokes, err := h.service.GetPokemons(r.Context())
	if err != nil {
		http.Error(w, "error obteniendo pokemones: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(pokes)
}

func (h *handler) CreatePokemon(w http.ResponseWriter, r *http.Request) {
	var pokemon domain.Pokemon
	if err := json.NewDecoder(r.Body).Decode(&pokemon); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	createdPokemon, err := h.service.Post(r.Context(), &pokemon)
	if err != nil {
		http.Error(w, "error creando pokemon", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(createdPokemon)
}

func (h *handler) GetPokemon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID no proporcionado", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	pokemon, err := h.service.GetPokemon(r.Context(), id)
	if err != nil {
		http.Error(w, "error obteniendo pokemon: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(pokemon)
}

func (h *handler) UpdatePokemon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID no proporcionado", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var pokemon domain.Pokemon
	if err := json.NewDecoder(r.Body).Decode(&pokemon); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = h.service.Patch(r.Context(), id, &pokemon)
	if err != nil {
		http.Error(w, "error actualizando pokemon", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) DeletePokemon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID no proporcionado", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, "error eliminando pokemon", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) GetTypes(w http.ResponseWriter, r *http.Request) {
	types, err := h.service.GetType(r.Context())
	if err != nil {
		http.Error(w, "error obteniendo tipos: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(types)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validación simple (hardcoded)
	if req.Username == "admin" && req.Password == "password" {
		token := "mi_secreto_token_seguro_12345" // Debe coincidir con .env
		resp := LoginResponse{Token: token}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
}
