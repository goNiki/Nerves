package env

import (
	"fmt"

	"github.com/caarlos0/env/v10"
	errorapp "github.com/goNiki/Nerves/internal/models/errorApp"
)

type apiEnvConfig struct {
	APIKey string `env:"API_KEY,required"`
}

type apiConfig struct {
	raw apiEnvConfig
}

func NewAPIConfig() (*apiConfig, error) {
	var raw apiEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("%w: %v", errorapp.ErrAPIConfig, err)
	}
	return &apiConfig{
		raw: raw,
	}, nil
}

func (cfg *apiConfig) APIKey() string {
	return cfg.raw.APIKey
}
