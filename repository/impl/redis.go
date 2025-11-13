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

func (r *Redis) AddBackList(userId string, token string, duration time.Duration) error {
	r.RedisClient.Set(context.Background(), userId, token, duration)
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

func (r *Redis) CreateRecord(userId uuid.UUID, accessToken string, accessTokenTimeLife time.Duration) error {
	status := r.RedisClient.Set(context.Background(), userId.String(), accessToken, accessTokenTimeLife)
	if err := status.Err(); err != nil {
		return fmt.Errorf("failed to create record: %v", err)
	}
	return nil
}
