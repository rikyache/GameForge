package service

import (
	"testsmth/internal/apperrors"
	"testsmth/internal/models"
)

type PlayerService struct {
	repo PlayerRepository
}

type PlayerRepository interface {
	GetPlayer(id int64) (*models.Player, error)
	AddPlayer(name string) error
	RemovePlayer(id int64) error
	ListPlayers() ([]models.Player, error)
}

func NewPlayerService(repo PlayerRepository) *PlayerService {
	return &PlayerService{
		repo: repo,
	}
}

func (s *PlayerService) GetPlayer(id int64) (*models.Player, error) {
	if id <= 0 {
		return nil, apperrors.ErrInvalidInput
	}

	return s.repo.GetPlayer(id)
}

func (s *PlayerService) AddPlayer(player models.Player) error {
	if player.Name == "" {
		return apperrors.ErrInvalidInput
	}

	if len(player.Name) > 100 {
		return apperrors.ErrInvalidInput
	}

	if player.Balance < 0 {
		return apperrors.ErrInvalidInput
	}

	return s.repo.AddPlayer(player.Name)
}

func (s *PlayerService) RemovePlayer(id int64) error {
	if id <= 0 {
		return apperrors.ErrInvalidInput
	}

	return s.repo.RemovePlayer(id)
}

func (s *PlayerService) ListPlayers() ([]models.Player, error) {
	return s.repo.ListPlayers()
}
