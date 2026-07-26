package main

import (
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

	http.HandleFunc("/players", handlers.GetPlayers)

	db := database.Connect()
	defer db.Close()

	err := repository.CreateTable(db)
	if err != nil {
		log.Fatal(err)
	}

	err = repository.AddPlayer(db, "Kirill", 100)
	if err != nil {
		log.Fatal(err)
	}

	player, err := repository.ListPlayer(db, 1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(player)

	http.ListenAndServe(":8080", nil)
}
