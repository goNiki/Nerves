package queue

import (
	"time"

	infraredis "github.com/goNiki/Nerves/internal/infrastructure/redis"
	"github.com/redis/go-redis/v9"
)

const (
	webhookQueueKey = "webhook:queue"
	dequeueTimeout  = 5 * time.Second
)

type Repository struct {
	redis *redis.Client
}

func NewQueueRepository(redis *infraredis.Client) *Repository {
	return &Repository{
		redis: redis.GetClient(),
	}
}
