package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/johnquangdev/oauth2/repository/interfaces"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	RedisClient *redis.Client
}

func NewRedis(RedisClient *redis.Client) interfaces.Redis {
	return &Redis{
		RedisClient: RedisClient,
	}
}

func (r *Redis) GetRefreshToken(userID string) (string, error) {
	return "", nil
}

func (r *Redis) DeleteRefreshToken(userID string) error {
	return nil
}

func (r *Redis) AddBackList(userID string, token string, duration time.Duration) error {
	r.RedisClient.Set(context.Background(), userID, token, duration)
	return nil
}

func (r *Redis) IsTokenBlacklisted(tokenID uuid.UUID) (bool, error) {
	key := "blacklist:" + tokenID.String()

	exists, err := r.RedisClient.Exists(context.Background(), key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check token blacklist: %w", err)
	}

	return exists == 1, nil
}
