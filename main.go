package main

import (
	"fmt"
	"log"
	"net/http"
	"testsmth/internal/cache"
	"testsmth/internal/database"
	"testsmth/internal/handlers"
	"testsmth/internal/middleware"
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
	//Подключаемся к postgresql
	db := database.Connect()
	defer db.Close()
	//Подключаемся к RedisClient
	redisClient := cache.Connect()
	defer redisClient.Close()

	redisCache := cache.NewRedisCache(redisClient)

	//base url

	playerRepo := repository.NewPlayerRepository(db)
	cachedPlayerRepo := repository.NewCachedPlayerRepository(playerRepo, redisCache)
	playerService := service.NewPlayerService(cachedPlayerRepo)
	playerHandler := handlers.NewPlayerHandler(playerService)

	gameRepo := repository.NewGameRepository(db)
	gameService := service.NewGameService(gameRepo)
	gameHandler := handlers.NewGameHandler(gameService)

	playerGamesRepo := repository.NewPlayerGameRepository(db, playerRepo, gameRepo)
	playerGamesService := service.NewPlayerGamesService(playerGamesRepo)
	playerGamesHandler := handlers.NewPlayerGamesHandler(playerGamesService)

	mux := http.NewServeMux()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello Docker")
	})

	//игроки
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello Docker")
	})

	// Players
	mux.HandleFunc("GET /players", playerHandler.ListPlayers)
	mux.HandleFunc("POST /players", playerHandler.AddPlayer)

	mux.HandleFunc("GET /players/{id}", playerHandler.GetOrDeletePlayer)
	mux.HandleFunc("DELETE /players/{id}", playerHandler.GetOrDeletePlayer)

	// Player actions
	mux.HandleFunc("POST /players/{id}/deposit", playerHandler.Deposit)

	// Games
	mux.HandleFunc("GET /games", gameHandler.ListGames)
	mux.HandleFunc("POST /games", gameHandler.AddGame)

	mux.HandleFunc("GET /games/{id}", gameHandler.GetOrDeleteGame)
	mux.HandleFunc("DELETE /games/{id}", gameHandler.GetOrDeleteGame)

	// Player-Games
	mux.HandleFunc("POST /players/{id}/games", playerGamesHandler.AddGame)
	mux.HandleFunc("GET /players/{id}/games", playerGamesHandler.GetPlayerGames)
	mux.HandleFunc("DELETE /players/{playerID}/games/{gameID}", playerGamesHandler.RemoveGame)
	mux.HandleFunc("POST /player/{playerID}/games/{gameID}/refund", playerGamesHandler.RefundGame)

	http.ListenAndServe(":8080", middleware.Logger(mux))
}
