package service

import (
	"errors"
	"testsmth/internal/models"
	"testsmth/internal/repository"
)

type PlayerService struct {
	repo *repository.PlayerRepository
}

func NewPlayerService(repo *repository.PlayerRepository) *PlayerService {
	return &PlayerService{
		repo: repo,
	}
}

func (s *PlayerService) GetPlayer(id int64) (*models.Player, error) {
	if id < 0 {
		return nil, ErrInvalidPlayerID
	}
	return s.repo.GetPlayer(id)
}

func (s *PlayerService) AddPlayer(player models.Player) error {
	if player.Name == "" {
		return errors.New("player name is empty")
	}
	if len(player.Name) > 100 {
		return ErrPlayerNameTooLong
	}
	return s.repo.AddPlayer(player.Name)
}

func (s *PlayerService) RemovePlayer(id int64) error {
	if id < 0 {
		return ErrInvalidPlayerID
	}
	return s.repo.RemovePlayer(id)
}

func (s *PlayerService) ListPlayers() ([]models.Player, error) {
	return s.repo.ListPlayers()
}
