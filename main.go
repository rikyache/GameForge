package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"testsmth/internal/database"
	"testsmth/internal/repository"

	_ "github.com/lib/pq"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello Docker")
	})

	db := database.Connect()
	defer db.Close()

	repository.AddPlayer(db, "Stas", 200)

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

	err := repository.CreateTable(db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(player)

	http.ListenAndServe(":8080", nil)
}
