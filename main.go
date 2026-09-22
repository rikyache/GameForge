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

	playerGamesRepo := repository.NewPlayerGameRepository(
		db,
		playerRepo,
		gameRepo,
	)

	cachedPlayerGamesRepo := repository.NewCachedPlayerGameRepository(
		playerGamesRepo,
		redisCache,
	)

	playerGamesService := service.NewPlayerGamesService(
		cachedPlayerGamesRepo,
	)

	playerGamesHandler := handlers.NewPlayerGamesHandler(
		playerGamesService,
	)

	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo)
	authHandler := handlers.NewAuthHandler(authService)

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
	mux.Handle("POST /players/{id}/deposit", middleware.Auth(http.HandlerFunc(playerHandler.Deposit)))
	mux.Handle("GET /players/{id}/profile", middleware.Auth(http.HandlerFunc(playerHandler.PlayerProfile)))

	// Games
	mux.HandleFunc("GET /games", gameHandler.ListGames)
	mux.Handle("POST /games", middleware.Auth(http.HandlerFunc(playerGamesHandler.AddGame)))

	mux.Handle("GET /games/{id}", middleware.Auth(http.HandlerFunc(gameHandler.GetOrDeleteGame)))
	mux.Handle("DELETE /games/{id}", middleware.Auth(http.HandlerFunc(gameHandler.GetOrDeleteGame)))

	// Player-Games
	mux.Handle("POST /players/{id}/games", middleware.Auth(http.HandlerFunc(playerGamesHandler.AddGame)))
	mux.Handle("GET /players/{id}/games", middleware.Auth(http.HandlerFunc(playerGamesHandler.GetPlayerGames)))
	mux.Handle("DELETE /players/{playerID}/games/{gameID}", middleware.Auth(http.HandlerFunc(playerGamesHandler.RemoveGame)))
	mux.Handle("POST /player/{playerID}/games/{gameID}/refund", middleware.Auth(http.HandlerFunc(playerGamesHandler.RefundGame)))

	// Register, login
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)

	http.ListenAndServe(":8080", middleware.Logger(mux))
}
