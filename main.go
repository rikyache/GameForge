package main

import (
	"fmt"
	"log"
	"net/http"
	"testsmth/internal/database"
	"testsmth/internal/handlers"

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
	playerHandler := handlers.PlayerHandler{
		DB: db,
	}

	gamesHandler := handlers.GameHandler{
		DB: db,
	}
	//игроки
	http.HandleFunc("/player/", playerHandler.GetOrDeletePlayer)
	http.HandleFunc("/players", playerHandler.ListPlayers)
	http.HandleFunc("/players/add", playerHandler.AddPlayer)
	//игры
	http.HandleFunc("/games", gamesHandler.ListGames)
	http.HandleFunc("/game/", gamesHandler.GetOrDeleteGame)
	http.HandleFunc("/games/add", gamesHandler.AddGame)

	http.ListenAndServe(":8080", nil)
}
