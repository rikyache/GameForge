package repository

import (
	"database/sql"
	"testing"
)

func TestAddPlayer(t *testing.T) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=kirill password=12345 dbname=practice sslmode=disable")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		t.Error(err)
	}
	tests := []struct {
		name       string
		playerName string
		wantErr    bool
	}{
		{
			name:       "add player",
			playerName: "Kirill",
			wantErr:    false,
		},
		{
			name:       "add player with empty name",
			playerName: "",
			wantErr:    false,
		},
	}

	playerRepo := NewPlayerRepository(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := playerRepo.AddPlayer(tt.playerName)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			var name string

			err = db.QueryRow(`
				SELECT name
				FROM players
				WHERE name = $1
			`, tt.playerName).Scan(&name)

			if err != nil {
				t.Fatal(err)
			}

			if name != tt.playerName {
				t.Errorf("got name %q, want %q", name, tt.playerName)
			}
		})
	}
}

func TestRemovePlayer(t *testing.T) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=kirill password=12345 dbname=practice sslmode=disable")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		t.Error(err)
	}

	tests := []struct {
		name        string
		playerExist bool
		wantErr     bool
		wantPlayer  bool
	}{
		{
			name:        "remove existing player",
			playerExist: true,
			wantErr:     false,
			wantPlayer:  false,
		},
		{
			name:        "remove non existing player",
			playerExist: false,
			wantErr:     false,
			wantPlayer:  false,
		},
	}

	playerRepo := NewPlayerRepository(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Arrange
			var playerID int64

			if tt.playerExist {
				err := db.QueryRow(`
					INSERT INTO players (name, balance)
					VALUES ($1, $2)
					RETURNING id
				`, "TestPlayer", 1000).Scan(&playerID)

				if err != nil {
					t.Fatal(err)
				}
			} else {
				playerID = 9999999
			}

			// Act
			err := playerRepo.RemovePlayer(playerID)

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

			// Assert: player exists
			var count int

			err = db.QueryRow(`
				SELECT COUNT(*)
				FROM players
				WHERE id = $1
			`, playerID).Scan(&count)

			if err != nil {
				t.Fatal(err)
			}

			if tt.wantPlayer && count != 1 {
				t.Errorf("player does not exist, but should")
			}

			if !tt.wantPlayer && count != 0 {
				t.Errorf("player still exists, but should be removed")
			}
		})
	}
}
