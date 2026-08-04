package service

import (
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

func (s *PlayerService) GetPlayer(id int) (*models.Player, error) {
	return s.repo.GetPlayer(id)
}

func (s *PlayerService) AddPlayer(player models.Player) error {
	return s.repo.AddPlayer(player.Name)
}

func (s *PlayerService) RemovePlayer(id int) error {
	return s.repo.RemovePlayer(id)
}

func (s *PlayerService) ListPlayers() ([]models.Player, error) {
	return s.repo.ListPlayers()
}
