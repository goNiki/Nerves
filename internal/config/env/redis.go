package env

import (
	"fmt"
	"net"
	"time"

	"github.com/caarlos0/env/v10"
	errorapp "github.com/goNiki/Nerves/internal/models/errorApp"
)

type RedisEnvConfig struct {
	Host            string        `env:"REDIS_HOST" envDefault:"localhost"`
	Port            string        `env:"REDIS_PORT" envDefault:"6379"`
	Password        string        `env:"REDIS_PASSWORD" envDefault:""`
	DB              int           `env:"REDIS_DB" envDefault:"0"`
	MaxRetries      int           `env:"REDIS_MAX_RETRIES" envDefault:"3"`
	PoolSize        int           `env:"REDIS_POOL_SIZE" envDefault:"10"`
	MinIdleConns    int           `env:"REDIS_MIN_IDLE_CONNS" envDefault:"5"`
	DialTimeout     time.Duration `env:"REDIS_DIAL_TIMEOUT" envDefault:"5s"`
	ReadTimeout     time.Duration `env:"REDIS_READ_TIMEOUT" envDefault:"3s"`
	WriteTimeout    time.Duration `env:"REDIS_WRITE_TIMEOUT" envDefault:"3s"`
	PoolTimeout     time.Duration `env:"REDIS_POOL_TIMEOUT" envDefault:"4s"`
	IdleTimeout     time.Duration `env:"REDIS_IDLE_TIMEOUT" envDefault:"5m"`
	MaxConnLifetime time.Duration `env:"REDIS_MAX_CONN_LIFETIME" envDefault:"30m"`
}

type RedisConfig struct {
	raw RedisEnvConfig
}

func NewRedisConfig() (*RedisConfig, error) {
	var raw RedisEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("%w: %v", errorapp.ErrParseRedisConfig, err)
	}

	return &RedisConfig{
		raw: raw,
	}, nil
}

func (cfg *RedisConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

func (cfg *RedisConfig) Password() string {
	return cfg.raw.Password
}

func (cfg *RedisConfig) DB() int {
	return cfg.raw.DB
}

func (cfg *RedisConfig) MaxRetries() int {
	return cfg.raw.MaxRetries
}

func (cfg *RedisConfig) PoolSize() int {
	return cfg.raw.PoolSize
}

func (cfg *RedisConfig) MinIdleConns() int {
	return cfg.raw.MinIdleConns
}

func (cfg *RedisConfig) DialTimeout() time.Duration {
	return cfg.raw.DialTimeout
}

func (cfg *RedisConfig) ReadTimeout() time.Duration {
	return cfg.raw.ReadTimeout
}

func (cfg *RedisConfig) WriteTimeout() time.Duration {
	return cfg.raw.WriteTimeout
}

func (cfg *RedisConfig) PoolTimeout() time.Duration {
	return cfg.raw.PoolTimeout
}

func (cfg *RedisConfig) IdleTimeout() time.Duration {
	return cfg.raw.IdleTimeout
}

func (cfg *RedisConfig) MaxConnLifetime() time.Duration {
	return cfg.raw.MaxConnLifetime
}
