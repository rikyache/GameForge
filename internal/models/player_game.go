package models

import "time"

type PlayerGame struct {
	PlayerID int64     `db:"player_id" json:"player_id"`
	GameID   int64     `db:"game_id" json:"game_id"`
	BoughtAt time.Time `db:"bought_at" json:"bought_at"`
}
