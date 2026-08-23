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

type CachedPlayerGameRepository struct {
	repo  *PlayerGameRepository
	cache cache.Cache
}

func NewCachedPlayerGameRepository(repo *PlayerGameRepository, cache cache.Cache) *CachedPlayerGameRepository {
	return &CachedPlayerGameRepository{
		repo:  repo,
		cache: cache,
	}
}

func (r *CachedPlayerGameRepository) BuyGame(playerID int64, gameID int64) error {
	if err := r.repo.BuyGame(playerID, gameID); err != nil {
		return err
	}

	ctx := context.Background()

	if err := r.cache.Delete(
		ctx,
		fmt.Sprintf("player:%d", playerID),
	); err != nil {
		return err
	}

	if err := r.cache.Delete(
		ctx,
		fmt.Sprintf("player_games:%d", playerID),
	); err != nil {
		return err
	}

	return nil
}

func (r *CachedPlayerGameRepository) GetPlayerGames(playerID int64) ([]models.OwnedGame, error) {
	key := fmt.Sprintf("playergames:%d", playerID)

	cached, err := r.cache.Get(context.Background(), key)
	if err != nil {
		return nil, err
	}
	if err == nil {
		var games []models.OwnedGame

		if err := json.Unmarshal([]byte(cached), &games); err != nil {
			return games, nil
		}

		log.Printf("cannot unmarshal cached games: %s", err)
	}

	if err != nil && err != cache.ErrCacheMiss {
		log.Printf("cannot get cached games from cache: %v", err)
	}

	games, err := r.repo.GetPlayerGames(playerID)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(games)
	if err != nil {
		return nil, err
	}

	if err := r.cache.Set(context.Background(), key, string(data), 5*time.Minute); err != nil {
		log.Printf("cannot set cached games: %s", err)
	}
	return games, nil
}

func (r *CachedPlayerGameRepository) Refund(
	playerID int64,
	gameID int64,
) error {

	if err := r.repo.Refund(playerID, gameID); err != nil {
		return err
	}

	ctx := context.Background()

	if err := r.cache.Delete(
		ctx,
		fmt.Sprintf("player:%d", playerID),
	); err != nil {
		return err
	}

	if err := r.cache.Delete(
		ctx,
		fmt.Sprintf("player_games:%d", playerID),
	); err != nil {
		return err
	}

	return nil
}
