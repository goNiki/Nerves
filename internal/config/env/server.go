package env

import (
	"fmt"
	"net"
	"time"

	"github.com/caarlos0/env/v10"
	errorapp "github.com/goNiki/Nerves/internal/models/errorApp"
)

type servercEnvConfig struct {
	Host            string        `env:"SERVER_HOST,required"`
	Port            string        `env:"SERVER_PORT,required"`
	ReadTimeOut     time.Duration `env:"SERVER_READ_TIMEOUT,required"`
	WriteTimeOut    time.Duration `env:"SERVER_WRITE_TIMEOUT,required"`
	IdleTimeOut     time.Duration `env:"SERVER_IDLE_TIMEOUT,required"`
	ShutDownTimeOut time.Duration `env:"SERVER_SHUTDOWN_TIMEOUT,required"`
}

type serverConfig struct {
	raw servercEnvConfig
}

func NewServerConfig() (*serverConfig, error) {
	var raw servercEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("%w: %v", errorapp.ErrServerConfig, err)
	}
	return &serverConfig{
		raw: raw,
	}, nil
}

func (cfg *serverConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

func (cfg *serverConfig) Port() string {
	return cfg.raw.Port
}

func (cfg *serverConfig) ReadTimeOut() time.Duration {
	return cfg.raw.ReadTimeOut
}

func (cfg *serverConfig) WriteTimeOut() time.Duration {
	return cfg.raw.WriteTimeOut
}

func (cfg *serverConfig) IdleTimeOut() time.Duration {
	return cfg.raw.IdleTimeOut
}

func (cfg *serverConfig) ShutDownTimeOut() time.Duration {
	return cfg.raw.ShutDownTimeOut
}
