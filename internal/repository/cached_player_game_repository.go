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

func NewCachedPlayerGameRepository(
	repo *PlayerGameRepository,
	cache cache.Cache,
) *CachedPlayerGameRepository {
	return &CachedPlayerGameRepository{
		repo:  repo,
		cache: cache,
	}
}

func (r *CachedPlayerGameRepository) BuyGame(ctx context.Context, playerID int64, gameID int64) error {
	if err := r.repo.BuyGame(ctx, playerID, gameID); err != nil {
		return err
	}

	if err := r.cache.Delete(ctx, fmt.Sprintf("player:%d", playerID)); err != nil {
		log.Printf("cannot invalidate player cache: %v", err)
	}

	if err := r.cache.Delete(ctx, fmt.Sprintf("player_games:%d", playerID)); err != nil {
		log.Printf("cannot invalidate player games cache: %v", err)
	}

	return nil
}

func (r *CachedPlayerGameRepository) GetPlayerGames(ctx context.Context, playerID int64) ([]models.OwnedGame, error) {
	key := fmt.Sprintf("player_games:%d", playerID)

	cached, err := r.cache.Get(ctx, key)

	if err == nil {
		var games []models.OwnedGame

		if err := json.Unmarshal([]byte(cached), &games); err == nil {
			return games, nil
		}

		log.Printf("cannot unmarshal cached games: %v", err)
	}

	if err != nil && err != cache.ErrCacheMiss {
		log.Printf("cannot get cached games: %v", err)
	}

	games, err := r.repo.GetPlayerGames(ctx, playerID)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(games)
	if err != nil {
		return nil, err
	}

	if err := r.cache.Set(ctx, key, string(data), 5*time.Minute); err != nil {
		log.Printf("cannot set cached games: %v", err)
	}

	return games, nil
}

func (r *CachedPlayerGameRepository) Refund(ctx context.Context, playerID int64, gameID int64) error {

	if err := r.repo.Refund(ctx, playerID, gameID); err != nil {
		return err
	}

	if err := r.cache.Delete(ctx, fmt.Sprintf("player:%d", playerID)); err != nil {
		log.Printf("cannot invalidate player cache: %v", err)
	}

	if err := r.cache.Delete(ctx, fmt.Sprintf("player_games:%d", playerID)); err != nil {
		log.Printf("cannot invalidate player games cache: %v", err)
	}

	return nil
}

func (r *CachedPlayerGameRepository) RemoveGame(
	ctx context.Context,
	playerID int64,
	gameID int64,
) error {

	if err := r.repo.RemoveGame(ctx, playerID, gameID); err != nil {
		return err
	}

	if err := r.cache.Delete(ctx, fmt.Sprintf("player_games:%d", playerID)); err != nil {
		log.Printf("cannot invalidate player games cache: %v", err)
	}

	return nil
}
