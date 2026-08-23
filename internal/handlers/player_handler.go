package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testsmth/internal/models"
	"testsmth/internal/service"
)

type PlayerHandler struct {
	Service *service.PlayerService
}

func NewPlayerHandler(service *service.PlayerService) *PlayerHandler {
	return &PlayerHandler{
		Service: service,
	}
}

func (h *PlayerHandler) GetOrDeletePlayer(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid player id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		player, err := h.Service.GetPlayer(id)
		if err != nil {
			handleError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(player); err != nil {
			handleError(w, err)
			return
		}

	case http.MethodDelete:
		err := h.Service.RemovePlayer(id)
		if err != nil {
			handleError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *PlayerHandler) ListPlayers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	players, err := h.Service.ListPlayers()
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(players); err != nil {
		handleError(w, err)
		return
	}
}

func (h *PlayerHandler) AddPlayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var player models.Player

	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.Service.AddPlayer(player)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *PlayerHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var amount models.DepositRequest

	err := json.NewDecoder(r.Body).Decode(&amount)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	money := amount.Amount
	idStr := r.PathValue("id")

	playerID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid player id", http.StatusBadRequest)
		return
	}

	err = h.Service.Deposit(playerID, money)
	if err != nil {
		handleError(w, err)
		return
	}
	return
}
