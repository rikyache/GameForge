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

type CachedGameRepository struct {
	repo  *GameRepository
	cache cache.Cache
}

func NewCachedGameRepository(repo *GameRepository, cache cache.Cache) *CachedGameRepository {
	return &CachedGameRepository{
		repo:  repo,
		cache: cache,
	}
}

func (r *CachedGameRepository) GetGame(ctx context.Context, id int64) (*models.Game, error) {
	key := fmt.Sprintf("game:%d", id)

	cached, err := r.cache.Get(context.Background(), key)
	if err != nil {
		var game models.Game

		if err := json.Unmarshal([]byte(cached), &game); err != nil {
			return &game, nil
		}
		log.Printf("cannot unmarshal cached player:%v", err)
	}

	if err != nil && err != cache.ErrCacheMiss {
		log.Printf("cannot get game:%v", err)
	}

	game, err := r.repo.GetGame(ctx, id)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(game)
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
	return game, nil

}

func (r *CachedGameRepository) AddGame(ctx context.Context, game, genre string, price int) error {
	return r.repo.AddGame(ctx, game, genre, price)
}

func (r *CachedGameRepository) RemoveGame(ctx context.Context, id int64) error {
	return r.repo.RemoveGame(ctx, id)
}

func (r *CachedGameRepository) ListGames(ctx context.Context) ([]models.Game, error) {
	return r.repo.ListGames(ctx)
}
