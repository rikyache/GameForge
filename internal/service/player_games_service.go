package service

import (
	"testsmth/internal/apperrors"
	"testsmth/internal/models"
	"testsmth/internal/repository"
)

type PlayerGamesService struct {
	repo *repository.PlayerGameRepository
}

func NewPlayerGamesService(repo *repository.PlayerGameRepository) *PlayerGamesService {
	return &PlayerGamesService{
		repo: repo,
	}
}

func (s *PlayerGamesService) BuyGame(playerID, gameID int64) error {
	if playerID <= 0 {
		return apperrors.ErrInvalidInput
	}

	if gameID <= 0 {
		return apperrors.ErrInvalidInput
	}

	return s.repo.BuyGame(playerID, gameID)
}

func (s *PlayerGamesService) GetPlayerGames(playerID int64) ([]models.OwnedGame, error) {
	if playerID <= 0 {
		return nil, apperrors.ErrInvalidInput
	}

	return s.repo.GetPlayerGames(playerID)
}

func (s *PlayerGamesService) RemoveGame(playerID, gameID int64) error {
	if playerID <= 0 {
		return apperrors.ErrInvalidInput
	}

	if gameID <= 0 {
		return apperrors.ErrInvalidInput
	}

	return s.repo.RemoveGame(playerID, gameID)
}
