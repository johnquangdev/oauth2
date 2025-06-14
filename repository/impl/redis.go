package impl

import (
	"context"
	"time"

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

func (r *Redis) SaveRefreshToken(userID string, token string, duration time.Duration) error {
	r.RedisClient.Set(context.Background(), userID, token, duration)
	return nil
}
func (r *Redis) GetRefreshToken(userID string) (string, error) {
	return "", nil
}
func (r *Redis) DeleteRefreshToken(userID string) error {
	return nil
}
