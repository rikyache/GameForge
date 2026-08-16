package repository

import (
	"testing"

	_ "github.com/lib/pq"
)

func TestBuyGame(t *testing.T) {
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
			db := setupTestDB(t)

			playerRepo := NewPlayerRepository(db)
			gameRepo := NewGameRepository(db)

			playerGameRepo := NewPlayerGameRepository(
				db,
				playerRepo,
				gameRepo,
			)

			// Arrange

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
				err := db.QueryRow(`
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

			// Act

			err := playerGameRepo.BuyGame(playerID, gameID)

			// Assert: error

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			// Assert: balance

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

			// Assert: game ownership

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
				t.Errorf("game was not added to player_games")
			}

			if !tt.wantGame && count != 0 {
				t.Errorf("game was added to player_games, but should not be")
			}
		})
	}
}

func TestGetPlayerGames(t *testing.T) {
	tests := []struct {
		name      string
		gameCount int
		wantCount int
	}{
		{
			name:      "player has one game",
			gameCount: 1,
			wantCount: 1,
		},
		{
			name:      "player has many games",
			gameCount: 5,
			wantCount: 5,
		},
		{
			name:      "player does not have games",
			gameCount: 0,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)

			playerGameRepo := NewPlayerGameRepository(db, nil, nil)

			// Arrange

			var playerID int64

			err := db.QueryRow(`
				INSERT INTO players (name, balance)
				VALUES ($1, $2)
				RETURNING id
			`, "Tester", 1000).Scan(&playerID)

			if err != nil {
				t.Fatal(err)
			}

			for i := 0; i < tt.gameCount; i++ {
				var gameID int64

				err := db.QueryRow(`
					INSERT INTO games (name, genre, price)
					VALUES ($1, $2, $3)
					RETURNING id
				`, "TestGame", "Test", 500).Scan(&gameID)

				if err != nil {
					t.Fatal(err)
				}

				_, err = db.Exec(`
					INSERT INTO player_games (player_id, game_id)
					VALUES ($1, $2)
				`, playerID, gameID)

				if err != nil {
					t.Fatal(err)
				}
			}

			// Act

			games, err := playerGameRepo.GetPlayerGames(playerID)

			if err != nil {
				t.Fatal(err)
			}

			// Assert

			if len(games) != tt.wantCount {
				t.Errorf(
					"got %d games, want %d",
					len(games),
					tt.wantCount,
				)
			}
		})
	}
}

func TestRemoveGame(t *testing.T) {
	tests := []struct {
		name       string
		gameExists bool
		wantErr    bool
		wantGame   bool
	}{
		{
			name:       "game exists",
			gameExists: true,
			wantErr:    false,
			wantGame:   false,
		},
		{
			name:       "game does not exist",
			gameExists: false,
			wantErr:    true,
			wantGame:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)

			playerGameRepo := NewPlayerGameRepository(db, nil, nil)

			// Arrange

			var playerID int64

			err := db.QueryRow(`
				INSERT INTO players (name, balance)
				VALUES ($1, $2)
				RETURNING id
			`, "Tester", 1000).Scan(&playerID)

			if err != nil {
				t.Fatal(err)
			}

			var gameID int64

			err = db.QueryRow(`
				INSERT INTO games (name, genre, price)
				VALUES ($1, $2, $3)
				RETURNING id
			`, "TestGame", "Test", 100).Scan(&gameID)

			if err != nil {
				t.Fatal(err)
			}

			if tt.gameExists {
				_, err = db.Exec(`
					INSERT INTO player_games (player_id, game_id)
					VALUES ($1, $2)
				`, playerID, gameID)

				if err != nil {
					t.Fatal(err)
				}
			}

			// Act

			err = playerGameRepo.RemoveGame(playerID, gameID)

			// Assert: error

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			// Assert: game exists

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

			if count != 0 {
				t.Errorf("game was not removed")
			}
		})
	}
}

func TestAddGame(t *testing.T) {
	tests := []struct {
		name     string
		wantErr  bool
		wantGame bool
	}{
		{
			name:     "add game",
			wantErr:  false,
			wantGame: true,
		},
		{
			name:     "game already added",
			wantErr:  true,
			wantGame: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)

			playerGameRepo := NewPlayerGameRepository(db, nil, nil)

			// Arrange

			var playerID int64

			err := db.QueryRow(`
				INSERT INTO players (name, balance)
				VALUES ($1, $2)
				RETURNING id
			`, "Tester", 1000).Scan(&playerID)

			if err != nil {
				t.Fatal(err)
			}

			var gameID int64

			err = db.QueryRow(`
				INSERT INTO games (name, genre, price)
				VALUES ($1, $2, $3)
				RETURNING id
			`, "TestGame", "TestGenre", 100).Scan(&gameID)

			if err != nil {
				t.Fatal(err)
			}

			if tt.wantErr {
				_, err = db.Exec(`
					INSERT INTO player_games (player_id, game_id)
					VALUES ($1, $2)
				`, playerID, gameID)

				if err != nil {
					t.Fatal(err)
				}
			}

			// Act

			err = playerGameRepo.AddGame(playerID, gameID)

			// Assert: error

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			// Assert: game exists

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
				t.Errorf("game was not added")
			}

			if !tt.wantGame && count != 0 {
				t.Errorf("game was added, but should not be")
			}
		})
	}
}

func TestExists(t *testing.T) {
	tests := []struct {
		name      string
		gameAdded bool
		want      bool
	}{
		{
			name:      "game exists",
			gameAdded: true,
			want:      true,
		},
		{
			name:      "game does not exist",
			gameAdded: false,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)

			playerGameRepo := NewPlayerGameRepository(db, nil, nil)

			// Arrange

			var playerID int64

			err := db.QueryRow(`
				INSERT INTO players (name, balance)
				VALUES ($1, $2)
				RETURNING id
			`, "Tester", 1000).Scan(&playerID)

			if err != nil {
				t.Fatal(err)
			}

			var gameID int64

			err = db.QueryRow(`
				INSERT INTO games (name, genre, price)
				VALUES ($1, $2, $3)
				RETURNING id
			`, "TestGame", "TestGenre", 100).Scan(&gameID)

			if err != nil {
				t.Fatal(err)
			}

			// Если по условию игра должна существовать
			if tt.gameAdded {
				_, err = db.Exec(`
					INSERT INTO player_games (player_id, game_id)
					VALUES ($1, $2)
				`, playerID, gameID)

				if err != nil {
					t.Fatal(err)
				}
			}

			// Act

			exists, err := playerGameRepo.Exists(playerID, gameID)

			// Assert

			if err != nil {
				t.Fatal(err)
			}

			if exists != tt.want {
				t.Errorf("exists = %v, want %v", exists, tt.want)
			}
		})
	}
}
