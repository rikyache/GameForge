package models

import "time"

type Game struct {
	ID    int64  `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Genre string `db:"genre" json:"genre"`
	Price int    `db:"price" json:"price"`
}

type OwnedGame struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Genre    string    `json:"genre"`
	BoughtAt time.Time `json:"bought_at"`
}
