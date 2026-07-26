package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	http.HandleFunc("/players", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

		player, err := repository.GetPlayer(db, 1)
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

	err = repository.AddPlayer(db, "Kirill", 100)
	if err != nil {
		log.Fatal(err)
	}

	player, err := repository.GetPlayer(db, 1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(player)

	http.ListenAndServe(":8080", nil)
}
