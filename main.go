package main

import (
	"fmt"
	"log"
	"net/http"
	"testsmth/internal/database"
	"testsmth/internal/handlers"
	"testsmth/internal/repository"
	"testsmth/internal/service"

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

	//base url
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello Docker")
	})

	playerRepo := repository.NewPlayerRepository(db)
	playerService := service.NewPlayerService(playerRepo)
	playerHandler := handlers.NewPlayerHandler(playerService)

	gameRepo := repository.NewGameRepository(db)
	gameService := service.NewGameService(gameRepo)
	gameHandler := handlers.NewGameHandler(gameService)

	//игроки
	http.HandleFunc("/player/", playerHandler.GetOrDeletePlayer)
	http.HandleFunc("/players", playerHandler.ListPlayers)
	http.HandleFunc("/players/add", playerHandler.AddPlayer)
	//игры
	http.HandleFunc("/games", gameHandler.ListGames)
	http.HandleFunc("/game/", gameHandler.GetOrDeleteGame)
	http.HandleFunc("/games/add", gameHandler.AddGame)

	http.ListenAndServe(":8080", nil)
}
