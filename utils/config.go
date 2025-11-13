package utils

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// Database configuration
	Dsn      string `envconfig:"DSN"`
	Port     string `envconfig:"DB_PORT"`
	Host     string `envconfig:"DB_HOST"`
	Name     string `envconfig:"DB_NAME"`
	Password string `envconfig:"DB_PASSWORD"`

	// jwt configuration
	SecretKey            string `envconfig:"SECRET_KEY"`
	AccessTokenTimeLife  uint16 `envconfig:"ACCESS_TOKEN_TIME_LIFE"`
	RefreshTokenTimeLife uint16 `envconfig:"REFRESH_TOKEN_TIME_LIFE"`

	// Redis configuration
	RedisAddr     string `envconfig:"REDIS_ADDR"`
	RedisPassword string `envconfig:"REDIS_PASSWORD"`
	RedisPort     string `envconfig:"REDIS_PORT"`

	// OAuth2 Google configuration
	RedirectUrl_Google  string `envconfig:"REDIRECT_URL_GOOGLE"`
	ClientId_Google     string `envconfig:"CLIENT_ID_GOOGLE"`
	ClientSecret_Google string `envconfig:"CLIENT_SECRET_GOOGLE"`
	Scopes_Google       string `envconfig:"SCOPES_GOOGLE"`

	RedirectUrl_GitHub  string `envconfig:"REDIRECT_URL_GITHUB"`
	ClientId_GitHub     string `envconfig:"CLIENT_ID_GITHUB"`
	ClientSecret_GitHub string `envconfig:"CLIENT_SECRET_GITHUB"`
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
