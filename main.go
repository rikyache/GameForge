package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"testsmth/internal/database"
	"testsmth/internal/models"
	"testsmth/internal/repository"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	//initialization
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	db := database.Connect()
	defer db.Close()

	//handlers

	//base url
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello Docker")
	})

	//функция для получения или удаления человека по id
	http.HandleFunc("/player/", func(w http.ResponseWriter, r *http.Request) {

		//получение чистого id
		idStr := strings.TrimPrefix(r.URL.Path, "/player/")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case "GET":
			player, err := repository.GetPlayer(db, id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(player)

		case "DELETE":
			err := repository.RemovePlayer(db, id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		default:
			http.Error(w, "invalid method", http.StatusMethodNotAllowed)
			return
		}

	})

	//список всех людей
	http.HandleFunc("/players", func(w http.ResponseWriter, r *http.Request) {
		var player []models.Player
		player, err := repository.ListPlayers(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(player)

	})

	http.HandleFunc("/players/add", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "POST" {
			http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		}

		var player models.Player

		err := json.NewDecoder(r.Body).Decode(&player)
		if err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		err = repository.AddPlayer(db, player.Name, player.Health)
		w.WriteHeader(http.StatusCreated)
	})

	err = repository.CreateTable(db)
	if err != nil {
		log.Fatal(err)
	}

	http.ListenAndServe(":8080", nil)
}
