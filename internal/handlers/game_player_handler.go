package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testsmth/internal/middleware"
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

	ctx := r.Context()

	playerID, ok := middleware.PlayerIDFromContext(ctx)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var purchaseRequest models.PurchaseRequest

	err := json.NewDecoder(r.Body).Decode(&purchaseRequest)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.Service.BuyGame(ctx,
		playerID,
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

	ctx := r.Context()

	playerID, ok := middleware.PlayerIDFromContext(ctx)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	games, err := h.Service.GetPlayerGames(ctx, playerID)
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

	ctx := r.Context()

	playerID, ok := middleware.PlayerIDFromContext(ctx)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	gameID, err := strconv.ParseInt(r.PathValue("gameID"), 10, 64)
	if err != nil {
		http.Error(w, "invalid game id", http.StatusBadRequest)
		return
	}

	if err := h.Service.RemoveGame(ctx, playerID, gameID); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PlayerGamesHandler) RefundGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	ctx := r.Context()

	playerID, ok := middleware.PlayerIDFromContext(ctx)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	gameIDstr := r.PathValue("gameID")

	gameID, err := strconv.ParseInt(gameIDstr, 10, 64)
	if err != nil {
		http.Error(w, "invalid game id", http.StatusBadRequest)
		return
	}

	err = h.Service.Refund(ctx, playerID, gameID)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
