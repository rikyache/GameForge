package auth

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenBlacklist interface {
	Revoke(ctx context.Context, tokenID string, ttl time.Duration) error
	IsRevoked(ctx context.Context, tokenID string) (bool, error)
}

type RedisTokenBlacklist struct {
	client *redis.Client
}

func NewRedisTokenBlacklist(client *redis.Client) *RedisTokenBlacklist {
	return &RedisTokenBlacklist{
		client: client,
	}
}

func (b *RedisTokenBlacklist) Revoke(ctx context.Context, tokenID string, ttl time.Duration) error {
	key := "revoked_token:" + tokenID

	return b.client.Set(ctx, key, "", ttl).Err()
}

func (b *RedisTokenBlacklist) IsRevoked(ctx context.Context, tokenID string) (bool, error) {
	key := "revoked_token:" + tokenID

	count, err := b.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
