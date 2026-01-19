package env

import (
	"fmt"
	"net"
	"time"

	"github.com/caarlos0/env/v10"
	errorapp "github.com/goNiki/Nerves/internal/models/errorApp"
)

type webhookEnvConfig struct {
	Host       string        `env:"WEBHOOK_HOST,required"`
	Port       string        `env:"WEBHOOK_PORT,required"`
	Timeout    time.Duration `env:"WEBHOOK_TIMEOUT,required"`
	MaxRetries int           `env:"WEBHOOK_MAX_RETRIES,required"`
}

type webhookConfig struct {
	raw webhookEnvConfig
}

func NewWebHookConfig() (*webhookConfig, error) {
	var raw webhookEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("%w: %v", errorapp.ErrWebHookConfig, err)
	}

	return &webhookConfig{
		raw: raw,
	}, nil
}

func (w *webhookConfig) Address() string {
	return net.JoinHostPort(w.raw.Host, w.raw.Port)
}

func (w *webhookConfig) Host() string {
	return w.raw.Host
}

func (w *webhookConfig) Port() string {
	return w.raw.Port
}

func (w *webhookConfig) Timeout() time.Duration {
	return w.raw.Timeout
}

func (w *webhookConfig) MaxRetries() int {
	return w.raw.MaxRetries
}
