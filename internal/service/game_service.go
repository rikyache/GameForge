package service

import (
	"context"
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

func (s *GameService) AddGame(ctx context.Context, game models.Game) error {
	if game.Name == "" {
		return apperrors.ErrInvalidInput
	}

	if game.Genre == "" {
		return apperrors.ErrInvalidInput
	}

	if len(game.Name) > 255 {
		return apperrors.ErrInvalidInput
	}

	return s.repo.AddGame(ctx,
		game.Name,
		game.Genre,
		game.Price,
	)
}

func (s *GameService) DeleteGame(ctx context.Context, id int64) error {
	if id <= 0 {
		return apperrors.ErrInvalidInput
	}

	return s.repo.RemoveGame(ctx, id)
}

func (s *GameService) GetGame(ctx context.Context, id int64) (*models.Game, error) {
	if id <= 0 {
		return nil, apperrors.ErrInvalidInput
	}

	return s.repo.GetGame(ctx, id)
}

func (s *GameService) ListGames(ctx context.Context) ([]models.Game, error) {
	return s.repo.ListGames(ctx)
}
