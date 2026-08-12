package repository

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func TestBuyGame(t *testing.T) {
	db, err := sql.Open(
		"postgres",
		"host=localhost port=5432 user=kirill password=12345 dbname=practice sslmode=disable",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	playerRepo := NewPlayerRepository(db)
	gameRepo := NewGameRepository(db)

	playerGameRepo := NewPlayerGameRepository(
		db,
		playerRepo,
		gameRepo,
	)

	tests := []struct {
		name          string
		playerBalance int
		gamePrice     int
		wantBalance   int
		wantErr       bool
		wantGame      bool
		playerExists  bool
		gameExists    bool
		checkBalance  bool
	}{
		{
			name:          "successful purchase",
			playerBalance: 1000,
			gamePrice:     300,
			wantBalance:   700,
			wantErr:       false,
			playerExists:  true,
			gameExists:    true,
			checkBalance:  true,
			wantGame:      true,
		},
		{
			name:          "not enough money",
			playerBalance: 200,
			gamePrice:     300,
			wantBalance:   200,
			wantErr:       true,
			playerExists:  true,
			gameExists:    true,
			checkBalance:  true,
			wantGame:      false,
		},
		{
			name:          "exact match money",
			playerBalance: 200,
			gamePrice:     200,
			wantBalance:   0,
			wantErr:       false,
			playerExists:  true,
			gameExists:    true,
			checkBalance:  true,
			wantGame:      true,
		},
		{
			name:          "player does not exist",
			playerBalance: 1000,
			gamePrice:     100,
			wantBalance:   0,
			wantErr:       true,
			playerExists:  false,
			gameExists:    true,
			checkBalance:  false,
			wantGame:      false,
		},
		{
			name:          "game does not exist",
			playerBalance: 1000,
			gamePrice:     100,
			wantBalance:   1000,
			wantErr:       true,
			playerExists:  true,
			gameExists:    false,
			checkBalance:  true,
			wantGame:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var playerID int64

			if tt.playerExists {
				err := db.QueryRow(`
                INSERT INTO players (name, balance)
                VALUES ($1, $2)
                RETURNING id
            `, "TestPlayer", tt.playerBalance).Scan(&playerID)

				if err != nil {
					t.Fatal(err)
				}
			} else {
				playerID = 9999999
			}

			var gameID int64
			if tt.gameExists {
				err = db.QueryRow(`
                INSERT INTO games (name, genre, price)
                VALUES ($1, $2, $3)
                RETURNING id
            `, "TestGame", "test", tt.gamePrice).Scan(&gameID)

				if err != nil {
					t.Fatal(err)
				}
			} else {
				gameID = 9999999
			}

			err = playerGameRepo.BuyGame(playerID, gameID)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
			if tt.checkBalance {
				var balance int

				err = db.QueryRow(`
                SELECT balance
                FROM players
                WHERE id = $1
            `, playerID).Scan(&balance)

				if err != nil {
					t.Fatal(err)
				}

				if balance != tt.wantBalance {
					t.Errorf(
						"balance = %d, want %d",
						balance,
						tt.wantBalance,
					)
				}
			}
			var count int

			err = db.QueryRow(`
				SELECT COUNT(*)
				FROM player_games
				WHERE player_id = $1
				AND game_id = $2
			`, playerID, gameID).Scan(&count)

			if err != nil {
				t.Fatal(err)
			}

			if tt.wantGame && count != 1 {
				t.Errorf("game wsa not added to player_games")
			}
			if !tt.wantGame && count != 0 {
				t.Errorf("game was added to player_games, but should not be")
			}
		})
	}
}
