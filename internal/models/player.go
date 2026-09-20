package models

import "time"

type Player struct {
	ID        int64     `db:"id"		json:"id"`
	Name      string    `db:"name"	json:"name"`
	Balance   int       `db:"balance"	json:"balance"`
	CreatedAt time.Time `db:"created_at"	json:"created_at"`
}

type PlayerProfile struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Balance      int    `json:"balance"`
	GamesCount   int    `json:"games_count"`
	LibraryValue int    `json:"library_value"`
}

type DepositRequest struct {
	Amount int `json:"amount"`
}
