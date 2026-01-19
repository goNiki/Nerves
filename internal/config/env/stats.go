package env

import (
	"fmt"

	"github.com/caarlos0/env/v10"
	errorapp "github.com/goNiki/Nerves/internal/models/errorApp"
)

type statscEnvConfig struct {
	TimeWindiwMinutes int `env:"STATS_TIME_WINDOW_MINUTES,required"`
}

type statsConfig struct {
	raw statscEnvConfig
}

func NewStatsConfig() (*statsConfig, error) {
	var raw statscEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("%w: %v", errorapp.ErrStatsConfig, err)
	}

	return &statsConfig{
		raw: raw,
	}, nil
}

func (s *statsConfig) TimeWindiwMinutes() int {
	return s.raw.TimeWindiwMinutes
}
