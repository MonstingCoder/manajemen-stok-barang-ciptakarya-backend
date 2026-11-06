package config

import (
	"errors"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var ErrParseEnv = errors.New("ERROR PARSE ENV")

type Config struct {
	DBHost            *string        `env:"DB_HOST" envDefault:"localhost"`
	DBUser            *string        `env:"DB_USER"`
	DBPassword        *string        `env:"DB_PASSWORD"`
	DBName            *string        `env:"DB_NAME"`
	DBTimeZone        *string        `env:"DB_TIMEZONE" envDefault:"Asia/Jakarta"`
	DBMaxOpenCons     *uint          `env:"DB_MAX_OPEN_CONS" envDefault:"5"`
	DBConnMaxIdleTime *time.Duration `env:"DB_CONN_MAX_IDLE_TIME" envDefault:"10m"`

	JWTSecret     *string        `env:"JWT_SECRET"`
	JWTExpireTime *time.Duration `env:"JWT_EXPIRE_TIME" envDefault:"30m"`

	AppPort        *string  `env:"APP_PORT" envDefault:":8080"`
	AllowedOrigins []string `env:"ALLOWED_ORIGINS" envDefault:"http://localhost:5173" envSeparator:"|"`
	AllowedMethods *string  `env:"ALLOWED_METHODS" envDefault:"GET, POST, PUT, PATCH, DELETE, OPTIONS"`
	AllowedHeaders *string  `env:"ALLOWED_HEADERS" envDefault:"Content-Type, Authorization"`

	DefaultDuration *time.Duration `env:"DEFAULT_DURATION" envDefault:"5s"`
}

func NewConfig() (cfg *Config, err error) {
	_ = godotenv.Load()

	cfg = new(Config)
	if *cfg, err = env.ParseAs[Config](); err != nil {
		err = ErrParseEnv
	}

	return
}
