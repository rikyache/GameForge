package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"testsmth/internal/models"
	"testsmth/internal/repository"
)

type PlayerHandler struct {
	DB *sql.DB
}

func (h *PlayerHandler) GetOrDeletePlayer(w http.ResponseWriter, r *http.Request) {
	//получение чистого id
	idStr := strings.TrimPrefix(r.URL.Path, "/player/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		player, err := repository.GetPlayer(h.DB, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(player)

	case "DELETE":
		err := repository.RemovePlayer(h.DB, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
}

func (h *PlayerHandler) ListPlayers(w http.ResponseWriter, r *http.Request) {
	var player []models.Player
	player, err := repository.ListPlayers(h.DB)
	if err != nil {
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

	var player models.Player

	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err = repository.AddPlayer(h.DB, player.Name)
	w.WriteHeader(http.StatusCreated)
}
