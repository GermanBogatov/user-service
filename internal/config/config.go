package config

import (
	"errors"
	"github.com/ilyakaznacheev/cleanenv"
)

var Namespace = "user_service"

type Health struct {
	CheckIntervalSec int `env:"USER_SERVICE_HEALTH_CHECK_INTERVAL" env-required:"10"`
}

type Postgres struct {
	URL                   string `env:"USER_SERVICE_POSTGRES_URL" env-required:"true"`
	ConnMaxLifetimeMinute int    `env:"USER_SERVICE_POSTGRES_CONN_MAX_LIFETIME_MINUTE" env-default:"1"`
	MaxOpenConn           int    `env:"USER_SERVICE_POSTGRES_MAX_OPEN_CONN" env-default:"10"`
	ConnAttempts          int    `env:"USER_SERVICE_POSTGRES_CONN_ATTEMPTS" env-default:"10"`
	ConnTimeout           int    `env:"USER_SERVICE_POSTGRES_CONN_TIMEOUT_SEC" env-default:"1"`
}

type Redis struct {
	Host       string `env:"USER_SERVICE_REDIS_HOST" env-required:"true"`
	Port       string `env:"USER_SERVICE_REDIS_PORT" env-required:"true"`
	Password   string `env:"USER_SERVICE_REDIS_PASSWORD"`
	DB         int    `env:"USER_SERVICE_REDIS_DB" env-default:"0"`
	UserTTL    int    `env:"USER_SERVICE_REDIS_USER_TTL" env-default:"20"`
	RefreshTTL int    `env:"USER_SERVICE_REDIS_REFRESH_TTL" env-default:"300"`
}

type Http struct {
	Port         string `env:"USER_SERVICE_HTTP_PORT" env-required:"true"`
	WriteTimeout int    `env:"USER_SERVICE_HTTP_WRITE_TIMEOUT_SEC" env-default:"60"`
	ReadTimeout  int    `env:"USER_SERVICE_HTTP_READ_TIMEOUT_SEC" env-default:"60"`
}

type Config struct {
	Health             Health
	Postgres           Postgres
	Redis              Redis
	Http               Http
	ShutdownTimeoutSec int `env:"USER_SERVICE_SHUTDOWN_TIMEOUT_SEC" env-default:"5"`
	JwtTTL             int `env:"USER_SERVICE_JWT_TTL" env-default:"300"`
}

// NewEnvConfig - парсим конфиг
func NewEnvConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	err = validateENV(cfg)
	if err != nil {
		return nil, err

	}

	return cfg, nil
}

func NewEnvConfigFromFile(filepath string) (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadConfig(filepath, cfg)
	if err != nil {
		return nil, err
	}

	err = validateENV(cfg)
	if err != nil {
		return nil, err

	}

	return cfg, nil

}

func validateENV(config *Config) error {
	if len(config.Postgres.URL) == 0 {
		return errors.New("empty postgres.url")
	}

	if len(config.Http.Port) == 0 {
		return errors.New("empty http.Port")
	}

	if len(config.Redis.Host) == 0 {
		return errors.New("empty redis.Host")
	}
	if len(config.Redis.Port) == 0 {
		return errors.New("empty redis.Port")
	}

	return nil
}
