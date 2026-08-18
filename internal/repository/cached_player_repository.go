package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testsmth/internal/cache"
	"testsmth/internal/models"
	"time"
)

type CachedPlayerRepository struct {
	repo  *PlayerRepository
	cache cache.Cache
}

func NewCachedPlayerRepository(repo *PlayerRepository, cache cache.Cache) *CachedPlayerRepository {
	return &CachedPlayerRepository{
		repo:  repo,
		cache: cache,
	}
}

func (r *CachedPlayerRepository) GetPlayer(id int64) (*models.Player, error) {
	key := fmt.Sprintf("player:%d", id)

	cached, err := r.cache.Get(context.Background(), key)
	if err == nil {
		var player models.Player

		if err := json.Unmarshal([]byte(cached), &player); err == nil {
			return &player, nil
		}
		log.Printf("cannot unmarshal cached player: %v", err)
	}

	if err != nil && err != cache.ErrCacheMiss {
		log.Printf("cannot get player: %v", err)
	}

	player, err := r.repo.GetPlayer(id)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(player)
	if err != nil {
		return nil, err
	}
	if err := r.cache.Set(
		context.Background(),
		key,
		string(data),
		5*time.Minute,
	); err != nil {

	}
	return player, nil
}

func (r *CachedPlayerRepository) AddPlayer(name string) error {
	return r.repo.AddPlayer(name)
}

func (r *CachedPlayerRepository) RemovePlayer(id int64) error {
	return r.repo.RemovePlayer(id)
}

func (r *CachedPlayerRepository) ListPlayers() ([]models.Player, error) {
	return r.repo.ListPlayers()
}
