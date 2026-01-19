package env

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v10"
	errorapp "github.com/goNiki/Nerves/internal/models/errorApp"
)

type DataBaseEnvConfig struct {
	Host              string        `env:"DB_HOST" envDefault:"localhost"`
	Port              string        `env:"DB_PORT" envDefault:"5432"`
	User              string        `env:"DB_USER" envDefault:"postgres"`
	Password          string        `env:"DB_PASSWORD" envDefault:"postgres"`
	Name              string        `env:"DB_NAME" envDefault:"postgres"`
	SslMode           string        `env:"DB_SSLMODE" envDefault:"disable"`
	MaxConns          int32         `env:"DB_MAXCONNS" envDefault:"20"`
	MinConns          int32         `env:"DB_MINCONNS" envDefault:"5"`
	MaxConnLifeTime   time.Duration `env:"DB_MAXCONNLIFETIME" envDefault:"30m"`
	MaxConnIdleTime   time.Duration `env:"DB_MAXCONNIDLETIME" envDefault:"5m"`
	HealthCheckPeriod time.Duration `env:"DB_HEALTHCHECKPERIOD" envDefault:"1m"`
}

type DBConfig struct {
	raw DataBaseEnvConfig
}

func NewDBConfig() (*DBConfig, error) {
	var raw DataBaseEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("%w: %v", errorapp.ErrParseDBConfig, err)
	}

	return &DBConfig{
		raw: raw,
	}, nil
}

func (cfx *DBConfig) Address() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfx.raw.User, cfx.raw.Password, cfx.raw.Host, cfx.raw.Port, cfx.raw.Name, cfx.raw.SslMode)
}

func (cfx *DBConfig) MaxConns() int32 {
	return cfx.raw.MaxConns
}

func (cfg *DBConfig) MinConns() int32 {
	return cfg.raw.MinConns
}

func (cfg *DBConfig) MaxConnLifeTime() time.Duration {
	return cfg.raw.MaxConnLifeTime
}

func (cfg *DBConfig) MaxConnIdleTime() time.Duration {
	return cfg.raw.MaxConnIdleTime
}

func (cft *DBConfig) HealthCheckPeriod() time.Duration {
	return cft.raw.HealthCheckPeriod
}
