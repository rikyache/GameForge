package service

import (
	"context"
	"testsmth/internal/apperrors"
	"testsmth/internal/models"
	"testsmth/internal/repository"
)

type PlayerGamesService struct {
	repo *repository.CachedPlayerGameRepository
}

func NewPlayerGamesService(
	repo *repository.CachedPlayerGameRepository,
) *PlayerGamesService {
	return &PlayerGamesService{
		repo: repo,
	}
}

func (s *PlayerGamesService) BuyGame(ctx context.Context, playerID, gameID int64) error {
	if playerID <= 0 {
		return apperrors.ErrInvalidInput
	}

	if gameID <= 0 {
		return apperrors.ErrInvalidInput
	}

	return s.repo.BuyGame(ctx, playerID, gameID)
}

func (s *PlayerGamesService) GetPlayerGames(ctx context.Context, playerID int64) ([]models.OwnedGame, error) {
	if playerID <= 0 {
		return nil, apperrors.ErrInvalidInput
	}

	return s.repo.GetPlayerGames(ctx, playerID)
}

func (s *PlayerGamesService) RemoveGame(ctx context.Context, playerID, gameID int64) error {
	if playerID <= 0 {
		return apperrors.ErrInvalidInput
	}

	if gameID <= 0 {
		return apperrors.ErrInvalidInput
	}

	return s.repo.RemoveGame(ctx, playerID, gameID)
}

func (s *PlayerGamesService) Refund(ctx context.Context, playerID, gameID int64) error {
	if playerID <= 0 {
		return apperrors.ErrInvalidInput
	}
	if gameID <= 0 {
		return apperrors.ErrInvalidInput
	}

	return s.repo.Refund(ctx, playerID, gameID)
}
