package repository

import (
	"context"

	"github.com/johnquangdev/oauth2/repository/impl"
	"github.com/johnquangdev/oauth2/repository/interfaces"
	"github.com/johnquangdev/oauth2/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type repository struct {
	db      *gorm.DB
	dbRedis *redis.Client
}

func (r repository) Auth() interfaces.Auth {
	return impl.NewAuth(r.db)
}

func (r repository) Redis() interfaces.Redis {
	return impl.NewRedis(r.dbRedis)
}

func NewRepository(db *gorm.DB, dbRedis *redis.Client) interfaces.Repo {
	return &repository{
		db:      db,
		dbRedis: dbRedis,
	}
}
func ConnectPostgres(config utils.Config) (*gorm.DB, error) {
	dsn := config.Dsn
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func ConnectRedis(cfg utils.Config) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
	})
	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	return redisClient, nil
}
