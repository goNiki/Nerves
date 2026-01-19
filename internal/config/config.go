package config

import (
	"fmt"

	"github.com/goNiki/Nerves/internal/config/env"
	errorapp "github.com/goNiki/Nerves/internal/models/errorApp"
	"github.com/subosito/gotenv"
)

type Config struct {
	Server   Server
	Postgres Postgres
	Redis    Redis
	Log      Log
	Api      Api
	Webhook  WebHook
	Stats    Stats
	Cache    Cache
}

func Load(path string) (*Config, error) {

	if err := gotenv.Load(path); err != nil {
		return nil, fmt.Errorf("%w: %v", errorapp.ErrLoadEnvFile, err)
	}

	server, err := env.NewServerConfig()
	if err != nil {
		return nil, err
	}

	api, err := env.NewAPIConfig()
	if err != nil {
		return nil, err
	}

	postgres, err := env.NewDBConfig()
	if err != nil {
		return nil, err
	}

	redis, err := env.NewRedisConfig()
	if err != nil {
		return nil, err
	}

	log, err := env.NewLoggerConfig()
	if err != nil {
		return nil, err
	}

	webhook, err := env.NewWebHookConfig()
	if err != nil {
		return nil, err
	}

	stats, err := env.NewStatsConfig()
	if err != nil {
		return nil, err
	}

	cache, err := env.NewCacheConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Server:   server,
		Postgres: postgres,
		Redis:    redis,
		Log:      log,
		Api:      api,
		Webhook:  webhook,
		Stats:    stats,
		Cache:    cache,
	}, nil

}
