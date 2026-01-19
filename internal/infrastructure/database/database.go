package database

import (
	"context"
	"fmt"

	"github.com/goNiki/Nerves/internal/config"
	errorapp "github.com/goNiki/Nerves/internal/models/errorApp"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	DB *pgxpool.Pool
}

func InitDatabase(cfg config.Postgres) (*Postgres, error) {

	pgxConf, err := pgxpool.ParseConfig(cfg.Address())
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errorapp.ErrInitPostgres, err)
	}

	pgxConf.MinConns = cfg.MinConns()
	pgxConf.MaxConns = cfg.MaxConns()
	pgxConf.MaxConnIdleTime = cfg.MaxConnIdleTime()
	pgxConf.MaxConnLifetime = cfg.MaxConnLifeTime()
	pgxConf.HealthCheckPeriod = cfg.HealthCheckPeriod()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, pgxConf)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errorapp.ErrInitPostgres, err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%w: %v", errorapp.ErrInitPostgres, err)
	}

	return &Postgres{
		DB: pool,
	}, nil
}
