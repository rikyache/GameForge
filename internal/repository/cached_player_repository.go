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

func (r *CachedPlayerRepository) GetPlayer(ctx context.Context, id int64) (*models.Player, error) {
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

	player, err := r.repo.GetPlayer(ctx, id)
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

func (r *CachedPlayerRepository) AddPlayer(ctx context.Context, name string) error {
	return r.repo.AddPlayer(ctx, name)
}

func (r *CachedPlayerRepository) RemovePlayer(ctx context.Context, id int64) error {
	return r.repo.RemovePlayer(ctx, id)
}

func (r *CachedPlayerRepository) ListPlayers(ctx context.Context) ([]models.Player, error) {
	return r.repo.ListPlayers(ctx)
}

func (r *CachedPlayerRepository) Deposit(ctx context.Context, playerID int64, amount int) error {
	if err := r.repo.Deposit(ctx, playerID, amount); err != nil {
		return err
	}

	key := fmt.Sprintf("player:%d", playerID)

	if err := r.cache.Delete(context.Background(), key); err != nil {
		return err
	}

	return nil
}

func (r *CachedPlayerRepository) GetProfile(ctx context.Context, playerID int64) (*models.PlayerProfile, error) {
	return r.repo.GetProfile(ctx, playerID)
}
