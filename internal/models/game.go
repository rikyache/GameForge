package models

import "time"

type Game struct {
	ID    int    `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Genre string `db:"genre" json:"genre"`
}

type OwnedGame struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Genre    string    `json:"genre"`
	BoughtAt time.Time `json:"bought_at"`
}
