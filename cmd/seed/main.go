package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	playersCount := flag.Int(
		"players",
		100_000,
		"number of players",
	)

	gamesCount := flag.Int(
		"games",
		1_000,
		"number of games",
	)

	playerGamesCount := flag.Int(
		"player-games",
		1_000_000,
		"number of player-game relations",
	)

	flag.Parse()

	db, err := sql.Open(
		"postgres",
		"host=localhost port=5432 user=kirill password=12345 dbname=practice_test sslmode=disable",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("connected to practice_test")

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	if err := seedPlayers(db, *playersCount, r); err != nil {
		log.Fatal(err)
	}

	if err := seedGames(db, *gamesCount, r); err != nil {
		log.Fatal(err)
	}

	if err := seedPlayerGames(
		db,
		*playerGamesCount,
		*playersCount,
		*gamesCount,
		r,
	); err != nil {
		log.Fatal(err)
	}

	log.Println("seed completed successfully")
}

func seedPlayers(db *sql.DB, count int, r *rand.Rand) error {
	log.Printf("creating %d players...", count)

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO players (name, balance)
		VALUES ($1, $2)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := 1; i <= count; i++ {
		name := fmt.Sprintf("Player_%d", i)
		balance := r.Intn(10_000)

		if _, err := stmt.Exec(name, balance); err != nil {
			return fmt.Errorf("insert player %d: %w", i, err)
		}

		if i%10_000 == 0 {
			log.Printf("players: %d/%d", i, count)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Println("players created")

	return nil
}

func seedGames(db *sql.DB, count int, r *rand.Rand) error {
	log.Printf("creating %d games...", count)

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO games (name, genre, price)
		VALUES ($1, $2, $3)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	genres := []string{
		"Action",
		"RPG",
		"Strategy",
		"Adventure",
		"Simulator",
		"Indie",
	}

	for i := 1; i <= count; i++ {
		name := fmt.Sprintf("Game_%d", i)
		genre := genres[r.Intn(len(genres))]
		price := r.Intn(5_000) + 100

		if _, err := stmt.Exec(name, genre, price); err != nil {
			return fmt.Errorf("insert game %d: %w", i, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Println("games created")

	return nil
}

func seedPlayerGames(
	db *sql.DB,
	count int,
	playersCount int,
	gamesCount int,
	r *rand.Rand,
) error {
	log.Printf("creating %d player_games...", count)

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO player_games (
			player_id,
			game_id,
			bought_at
		)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := 1; i <= count; i++ {
		playerID := r.Intn(playersCount) + 1
		gameID := r.Intn(gamesCount) + 1

		boughtAt := time.Now().Add(
			-time.Duration(r.Intn(365*24)) * time.Hour,
		)

		if _, err := stmt.Exec(
			playerID,
			gameID,
			boughtAt,
		); err != nil {
			return fmt.Errorf("insert player_game %d: %w", i, err)
		}

		if i%100_000 == 0 {
			log.Printf("player_games: %d/%d", i, count)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Println("player_games created")

	return nil
}
