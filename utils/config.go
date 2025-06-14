package utils

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Dsn      string `envconfig:"DSN"`
	Port     string `envconfig:"DB_PORT"`
	Host     string `envconfig:"DB_HOST"`
	Name     string `envconfig:"DB_NAME"`
	Password string `envconfig:"DB_PASSWORD"`

	SecretKey            string `envconfig:"SECRET_KEY"`
	AccessTokenTimeLife  uint16 `envconfig:"ACCESS_TOKEN_TIME_LIFE"`
	RefreshTokenTimeLife uint16 `envconfig:"RESFRESH_TOKEN_TIME_LIFE"`

	RedisAddr     string `envconfig:"REDIS_ADDR"`
	RedisPassword string `envconfig:"REDIS_PASSWORD"`
	RedisPort     string `envconfig:"REDIS_PORT"`
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return &Config{}, err
	}
	cfg := new(Config)
	err = envconfig.Process("", cfg)
	if err != nil {
		return &Config{}, err
	}
	return cfg, nil
}
