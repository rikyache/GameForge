package handlers

import (
	"encoding/json"
	"errors"
	"log"
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

func (h *PlayerGamesHandler) BuyGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var purchaseRequest models.PurchaseRequest

	err := json.NewDecoder(r.Body).Decode(&purchaseRequest)
	defer r.Body.Close()

	if err != nil {
		log.Printf("failed to parse body: %s", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	playerID := purchaseRequest.PlayerID
	gameID := purchaseRequest.GameID

	err = h.Service.BuyGame(playerID, gameID)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidPlayerID):
			http.Error(w, "invalid player id", http.StatusBadRequest)

		case errors.Is(err, service.ErrInvalidGameID):
			http.Error(w, "invalid game id", http.StatusBadRequest)

		case errors.Is(err, service.ErrGameNotFound):
			http.Error(w, "game not found", http.StatusNotFound)

		case errors.Is(err, service.ErrAlreadyExists):
			http.Error(w, "game already owned", http.StatusConflict)

		default:
			log.Printf("failed to buy game: %s", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *PlayerGamesHandler) GetPlayerGames(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/players/")
	path = strings.TrimSuffix(path, "/games")

	playerID, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		log.Printf("failed to parse path: %s", err)
		http.Error(w, "invalid player id", http.StatusBadRequest)
		return
	}

	games, err := h.Service.GetPlayerGames(playerID)
	if err != nil {
		log.Printf("failed to get games: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(games)
}

func (h *PlayerGamesHandler) RemoveGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 5 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	playerIDs := parts[2]
	gameIDs := parts[len(parts)-1]

	playerID, err := strconv.ParseInt(playerIDs, 10, 64)
	if err != nil {
		log.Printf("failed to parse playerID: %s", err)
		http.Error(w, "invalid player id", http.StatusBadRequest)
		return
	}
	gameID, err := strconv.ParseInt(gameIDs, 10, 64)
	if err != nil {
		log.Printf("failed to parse gameID: %s", err)
		http.Error(w, "invalid game id", http.StatusBadRequest)
		return
	}

	err = h.Service.RemoveGame(playerID, gameID)
	if err != nil {
		log.Printf("failed to remove game: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
