package service

import (
	"testsmth/internal/models"
	"testsmth/internal/repository"
)

type PlayerGamesService struct {
	repo *repository.PlayerGameRepository
}

func NewPlayerGamesService(repo *repository.PlayerGameRepository) *PlayerGamesService {
	return &PlayerGamesService{repo: repo}
}

func (s *PlayerGamesService) BuyGame(playerID int64, gameID int64) error {
	if playerID <= 0 {
		return ErrInvalidPlayerID
	}
	if gameID <= 0 {
		return ErrInvalidGameID
	}

	exists, err := s.repo.Exists(
		playerID,
		gameID,
	)
	if err != nil {
		return err
	}

	if exists {
		return ErrAlreadyExists
	}

	return s.repo.BuyGame(
		playerID,
		gameID,
	)
}

func (s *PlayerGamesService) GetPlayerGames(playerID int64) ([]models.OwnedGame, error) {
	if playerID <= 0 {
		return nil, ErrInvalidPlayerID
	}

	return s.repo.GetPlayerGames(playerID)
}

func (s *PlayerGamesService) RemoveGame(playerID int64, gameID int64) error {
	if playerID <= 0 {
		return ErrInvalidPlayerID
	}
	if gameID <= 0 {
		return ErrInvalidGameID
	}

	return s.repo.RemoveGame(playerID, gameID)
}
