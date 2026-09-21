package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"testsmth/internal/models"
	"testsmth/internal/service"
)

type GameHandler struct {
	Service *service.GameService
}

func NewGameHandler(service *service.GameService) *GameHandler {
	return &GameHandler{
		Service: service,
	}
}

func (h *GameHandler) ListGames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	games, err := h.Service.ListGames(ctx)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(games); err != nil {
		handleError(w, err)
		return
	}
}

func (h *GameHandler) GetOrDeleteGame(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/games/")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid game id", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	switch r.Method {

	case http.MethodGet:
		game, err := h.Service.GetGame(ctx, id)
		if err != nil {
			handleError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(game); err != nil {
			handleError(w, err)
			return
		}

	case http.MethodDelete:
		err := h.Service.DeleteGame(ctx, id)
		if err != nil {
			handleError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *GameHandler) AddGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	ctx := r.Context()

	var game models.Game

	err := json.NewDecoder(r.Body).Decode(&game)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.Service.AddGame(ctx, game)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
