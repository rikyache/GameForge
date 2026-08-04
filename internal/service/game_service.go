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
	return s.repo.AddGame(game.Name, game.Genre)
}

func (s *GameService) DeleteGame(id int) error {
	return s.repo.RemoveGame(id)
}

func (s *GameService) GetGame(id int) (*models.Game, error) {
	return s.repo.GetGame(id)
}

func (s *GameService) ListGames() ([]models.Game, error) {
	return s.repo.ListGames()
}
