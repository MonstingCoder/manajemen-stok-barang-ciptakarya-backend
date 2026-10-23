package config

import (
	"errors"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var ErrParseEnv = errors.New("ERROR PARSE ENV")

type Config struct {
	DBUser     string `env:"DB_USER"`
	DBPassword string
	DBName     string
}

func NewConfig() (cfg *Config, err error) {
	_ = godotenv.Load()

	if cfg, err = env.ParseAs[*Config](); err != nil {
		err = ErrParseEnv
	}

	return
}
