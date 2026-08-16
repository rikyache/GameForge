package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"testsmth/internal/models"
	"testsmth/internal/service"
)

type PlayerGamesHandler struct {
	Service *service.PlayerGamesService
}

func NewPlayerGamesHandler(service *service.PlayerGamesService) *PlayerGamesHandler {
	return &PlayerGamesHandler{
		Service: service,
	}
}

func (h *PlayerGamesHandler) AddGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var purchaseRequest models.PurchaseRequest

	err := json.NewDecoder(r.Body).Decode(&purchaseRequest)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.Service.BuyGame(
		purchaseRequest.PlayerID,
		purchaseRequest.GameID,
	)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *PlayerGamesHandler) GetPlayerGames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/players/")
	path = strings.TrimSuffix(path, "/games")

	playerID, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "invalid player id", http.StatusBadRequest)
		return
	}

	games, err := h.Service.GetPlayerGames(playerID)
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

func (h *PlayerGamesHandler) RemoveGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 5 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	playerID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "invalid player id", http.StatusBadRequest)
		return
	}

	gameID, err := strconv.ParseInt(parts[len(parts)-1], 10, 64)
	if err != nil {
		http.Error(w, "invalid game id", http.StatusBadRequest)
		return
	}

	err = h.Service.RemoveGame(playerID, gameID)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
