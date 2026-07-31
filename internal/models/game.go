package models

type Game struct {
	ID    int    `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Genre string `db:"genre" json:"genre"`
}
