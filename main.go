package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"testsmth/internal/database"
	"testsmth/internal/repository"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello Docker")
	})

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	db := database.Connect()
	defer db.Close()

	//функция для получения человека по id через Http
	http.HandleFunc("/players", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id, err := strconv.Atoi(
			r.URL.Query().Get("id"),
		)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		player, err := repository.GetPlayer(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(player)
	})

	err = repository.CreateTable(db)
	if err != nil {
		log.Fatal(err)
	}

	http.ListenAndServe(":8080", nil)
}
