package configs

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

const (
	Namespace = "user_service"
)

type Postgres struct {
	URL string `env:"USER_SERVICE_POSTGRES_URL"`
}

type Config struct {
	Postgres Postgres
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
	return nil
}
