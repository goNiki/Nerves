package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/goNiki/Nerves/internal/config"
	errorapp "github.com/goNiki/Nerves/internal/models/errorApp"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.Client
}

func New(cfg config.Redis) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:            cfg.Address(),
		Password:        cfg.Password(),
		DB:              cfg.DB(),
		MaxRetries:      cfg.MaxRetries(),
		PoolSize:        cfg.PoolSize(),
		MinIdleConns:    cfg.MinIdleConns(),
		DialTimeout:     cfg.DialTimeout(),
		ReadTimeout:     cfg.ReadTimeout(),
		WriteTimeout:    cfg.WriteTimeout(),
		PoolTimeout:     cfg.PoolTimeout(),
		ConnMaxIdleTime: cfg.IdleTimeout(),
		ConnMaxLifetime: cfg.MaxConnLifetime(),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errorapp.ErrInitRedis, err)
	}

	return &Client{rdb: rdb}, nil
}

func (c *Client) GetClient() *redis.Client {
	return c.rdb
}
