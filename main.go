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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello Docker")
	})

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

	//игроки
	http.HandleFunc("/player/", playerHandler.GetOrDeletePlayer)
	http.HandleFunc("/players", playerHandler.ListPlayers)
	http.HandleFunc("/players/add", playerHandler.AddPlayer)
	//игры
	http.HandleFunc("/games", gameHandler.ListGames)
	http.HandleFunc("/game/", gameHandler.GetOrDeleteGame)
	http.HandleFunc("/games/add", gameHandler.AddGame)
	//связь игры-игроки
	http.HandleFunc("/players/games", playerGamesHandler.AddGame)

	http.HandleFunc("/players/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			playerGamesHandler.GetPlayerGames(w, r)

		case http.MethodDelete:
			playerGamesHandler.RemoveGame(w, r)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8080", middleware.Logger(http.DefaultServeMux))
}
