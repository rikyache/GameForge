package service

import "errors"

var (
	ErrPlayerNameEmpty       = errors.New("player name is empty")
	ErrPlayerNameTooLong     = errors.New("player name is too long")
	ErrPlayerBalanceNegative = errors.New("player balance is negative")

	ErrGameNameEmpty   = errors.New("game name is empty")
	ErrGameNameTooLong = errors.New("game name is too long")
	ErrGameGenreEmpty  = errors.New("game genre is empty")

	ErrInvalidGameID   = errors.New("invalid game id")
	ErrInvalidPlayerID = errors.New("invalid player id")

	ErrAlreadyExists = errors.New("already exists")
	ErrGameNotFound  = errors.New("game not found for player")
)
