package service

import (
	"testsmth/internal/apperrors"
	"testsmth/internal/models"
	"testsmth/internal/repository"
)

type GameService struct {
	repo *repository.GameRepository
}

func NewGameService(repo *repository.GameRepository) *GameService {
	return &GameService{
		repo: repo,
	}
}

func (s *GameService) AddGame(game models.Game) error {
	if game.Name == "" {
		return apperrors.ErrInvalidInput
	}

	if game.Genre == "" {
		return apperrors.ErrInvalidInput
	}

	if len(game.Name) > 255 {
		return apperrors.ErrInvalidInput
	}

	return s.repo.AddGame(
		game.Name,
		game.Genre,
		game.Price,
	)
}

func (s *GameService) DeleteGame(id int64) error {
	if id <= 0 {
		return apperrors.ErrInvalidInput
	}

	return s.repo.RemoveGame(id)
}

func (s *GameService) GetGame(id int64) (*models.Game, error) {
	if id <= 0 {
		return nil, apperrors.ErrInvalidInput
	}

	return s.repo.GetGame(id)
}

func (s *GameService) ListGames() ([]models.Game, error) {
	return s.repo.ListGames()
}
