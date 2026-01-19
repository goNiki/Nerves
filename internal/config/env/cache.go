package env

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v10"
	errorapp "github.com/goNiki/Nerves/internal/models/errorApp"
)

type cacheEnvConfig struct {
	RefreshInterval time.Duration `env:"CACHE_REFRESH_INTERVAL" envDefault:"3m"`
}

type cacheConfig struct {
	raw cacheEnvConfig
}

func NewCacheConfig() (*cacheConfig, error) {
	var raw cacheEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("%w: %v", errorapp.ErrCacheConfig, err)
	}

	return &cacheConfig{
		raw: raw,
	}, nil
}

func (c *cacheConfig) RefreshInterval() time.Duration {
	return c.raw.RefreshInterval
}
