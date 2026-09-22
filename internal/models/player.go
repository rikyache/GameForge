package models

import "time"

type Player struct {
	ID           int64     `db:"id"			json:"id"`
	Name         string    `db:"name"			json:"name"`
	Balance      int       `db:"balance"		json:"balance"`
	Email        string    `db:"email"			json:"email"`
	PasswordHash string    `db:"password_hash"	json:"password_hash"`
	CreatedAt    time.Time `db:"created_at"	json:"created_at"`
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

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserParams struct {
	Name         string
	Email        string
	PasswordHash string
}
