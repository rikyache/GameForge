package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	//получение чистого id
	idStr := strings.TrimPrefix(r.URL.Path, "/player/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		log.Printf("failed to parse id: %v", err)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		player, err := h.Service.GetPlayer(id)
		if err != nil {
			log.Printf("failed to get player: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(player)

	case "DELETE":
		err := h.Service.RemovePlayer(id)
		if err != nil {
			log.Printf("failed to remove player: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
}

func (h *PlayerHandler) ListPlayers(w http.ResponseWriter, r *http.Request) {
	var player []models.Player
	player, err := h.Service.ListPlayers()
	if err != nil {
		log.Printf("failed to list players: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(player)
}

func (h *PlayerHandler) AddPlayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
	}

	var player *models.Player

	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		log.Printf("failed to parse body: %v", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err = h.Service.AddPlayer(*player)
	w.WriteHeader(http.StatusCreated)
}
