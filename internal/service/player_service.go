package service

import (
	"context"
	"testsmth/internal/apperrors"
	"testsmth/internal/models"
)

type PlayerService struct {
	repo PlayerRepository
}

type PlayerRepository interface {
	GetPlayer(ctx context.Context, id int64) (*models.Player, error)
	AddPlayer(ctx context.Context, name string) error
	RemovePlayer(ctx context.Context, id int64) error
	ListPlayers(ctx context.Context) ([]models.Player, error)
	Deposit(ctx context.Context, playerID int64, amount int) error
	GetProfile(ctx context.Context, playerID int64) (*models.PlayerProfile, error)
}

func NewPlayerService(repo PlayerRepository) *PlayerService {
	return &PlayerService{
		repo: repo,
	}
}

func (s *PlayerService) GetPlayer(ctx context.Context, id int64) (*models.Player, error) {
	if id <= 0 {
		return nil, apperrors.ErrInvalidInput
	}

	return s.repo.GetPlayer(ctx, id)
}

func (s *PlayerService) AddPlayer(ctx context.Context, player models.Player) error {
	if player.Name == "" {
		return apperrors.ErrInvalidInput
	}

	if len(player.Name) > 100 {
		return apperrors.ErrInvalidInput
	}

	if player.Balance < 0 {
		return apperrors.ErrInvalidInput
	}

	return s.repo.AddPlayer(ctx, player.Name)
}

func (s *PlayerService) RemovePlayer(ctx context.Context, id int64) error {
	if id <= 0 {
		return apperrors.ErrInvalidInput
	}

	return s.repo.RemovePlayer(ctx, id)
}

func (s *PlayerService) ListPlayers(ctx context.Context) ([]models.Player, error) {
	return s.repo.ListPlayers(ctx)
}

func (s *PlayerService) Deposit(ctx context.Context, playerID int64, amount int) error {
	if amount <= 0 {
		return apperrors.ErrInvalidInput
	}
	if playerID <= 0 {
		return apperrors.ErrInvalidInput
	}
	return s.repo.Deposit(ctx, playerID, amount)
}

func (s *PlayerService) GetProfile(ctx context.Context, playerID int64) (*models.PlayerProfile, error) {
	if playerID <= 0 {
		return nil, apperrors.ErrInvalidInput
	}
	return s.repo.GetProfile(ctx, playerID)
}
