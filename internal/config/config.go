package config

import (
	"errors"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var ErrParseEnv = errors.New("ERROR PARSE ENV")

type Config struct {
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_NAME"`

	JWTSecret     []byte `env:"JWT_SECRET"`
	JWTExpireTime time.Duration

	AllowedOrigins []string `env:"ALLOWED_ORIGINS"`
	AllowedMethods string
	AllowedHeaders string

	DefaultDuration time.Duration
}

func NewConfig() (cfg *Config, err error) {
	_ = godotenv.Load()

	if cfg, err = env.ParseAs[*Config](); err != nil {
		err = ErrParseEnv
	}

	cfg.JWTExpireTime = 30 * time.Minute
	cfg.AllowedOrigins = []string{"http://localhost:5173"}
	cfg.AllowedMethods = "GET, POST, PUT, DELETE, OPTIONS"
	cfg.AllowedHeaders = "Content-Type, Authorization"
	cfg.DefaultDuration = 5 * time.Second

	return
}
