package service

import (
	"testsmth/internal/models"
	"testsmth/internal/repository"
)

type GameService struct {
	repo *repository.GameRepository
}

func NewGameService(repo *repository.GameRepository) *GameService {
	return &GameService{repo: repo}
}

func (s *GameService) AddGame(game models.Game) error {
	if game.Name == "" {
		return ErrGameNameEmpty
	}
	if game.Genre == "" {
		return ErrGameGenreEmpty
	}
	if len(game.Name) > 255 {
		return ErrGameNameTooLong
	}
	return s.repo.AddGame(
		game.Name,
		game.Genre,
	)
}

func (s *GameService) DeleteGame(id int64) error {
	if id <= 0 {
		return ErrInvalidGameID
	}
	return s.repo.RemoveGame(id)
}

func (s *GameService) GetGame(id int64) (*models.Game, error) {
	if id <= 0 {
		return nil, ErrInvalidGameID
	}
	return s.repo.GetGame(id)
}

func (s *GameService) ListGames() ([]models.Game, error) {
	return s.repo.ListGames()
}
