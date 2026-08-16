package repository

import (
	"testing"
)

func TestGetGame(t *testing.T) {
	tests := []struct {
		name     string
		wantErr  bool
		wantGame bool
	}{
		{
			name:     "get existing game",
			wantErr:  false,
			wantGame: true,
		},
		{
			name:     "get non existing game",
			wantErr:  true,
			wantGame: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)

			gameRepo := NewGameRepository(db)

			var gameID int64

			if tt.wantGame {
				err := db.QueryRow(`
					INSERT INTO games (name, genre)
					VALUES ($1, $2)
					RETURNING id
				`, "Minecraft", "Sandbox").Scan(&gameID)

				if err != nil {
					t.Fatal(err)
				}
			} else {
				gameID = -1
			}

			game, err := gameRepo.GetGame(gameID)

			if (err != nil) != tt.wantErr {
				t.Errorf(
					"GetGame() error = %v, wantErr = %v",
					err,
					tt.wantErr,
				)
			}

			if (game != nil) != tt.wantGame {
				t.Errorf(
					"GetGame() game = %v, wantGame = %v",
					game,
					tt.wantGame,
				)
			}

			if tt.wantGame {
				if game.ID != gameID {
					t.Errorf(
						"GetGame() ID = %v, want %v",
						game.ID,
						gameID,
					)
				}

				if game.Name != "Minecraft" {
					t.Errorf(
						"GetGame() Name = %v, want Minecraft",
						game.Name,
					)
				}

				if game.Genre != "Sandbox" {
					t.Errorf(
						"GetGame() Genre = %v, want Sandbox",
						game.Genre,
					)
				}
			}
		})
	}
}

func TestListGames(t *testing.T) {
	tests := []struct {
		name       string
		createGame bool
		wantGames  int
	}{
		{
			name:       "list existing games",
			createGame: true,
			wantGames:  1,
		},
		{
			name:       "list empty games",
			createGame: false,
			wantGames:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)

			gameRepo := NewGameRepository(db)

			if tt.createGame {
				_, err := db.Exec(`
					INSERT INTO games (name, genre)
					VALUES ($1, $2)
				`, "Minecraft", "Sandbox")

				if err != nil {
					t.Fatal(err)
				}
			}

			games, err := gameRepo.ListGames()

			if err != nil {
				t.Fatal(err)
			}

			if len(games) != tt.wantGames {
				t.Errorf(
					"ListGames() returned %d games, want %d",
					len(games),
					tt.wantGames,
				)
			}

			if tt.createGame {
				if games[0].Name != "Minecraft" {
					t.Errorf(
						"Name = %v, want Minecraft",
						games[0].Name,
					)
				}

				if games[0].Genre != "Sandbox" {
					t.Errorf(
						"Genre = %v, want Sandbox",
						games[0].Genre,
					)
				}
			}
		})
	}
}
