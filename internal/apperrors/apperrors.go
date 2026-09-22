package apperrors

import "errors"

var (
	ErrInvalidInput = errors.New("invalid input")

	ErrPlayerNotFound = errors.New("player not found")
	ErrGameNotFound   = errors.New("game not found")

	ErrPlayerGameNotFound = errors.New("player does not own game")

	ErrAlreadyExists       = errors.New("resource already exists")
	ErrInsufficientBalance = errors.New("insufficient balance")

	ErrInvalidCredentials = errors.New("invalid credentials")
)
