package cacheincident

import (
	"time"

	infraredis "github.com/goNiki/Nerves/internal/infrastructure/redis"
	"github.com/redis/go-redis/v9"
)

const (
	activeIncidentsHash  = "active_incidents"
	incidentsGeo         = "incidents:geo"
	maxIncidentRadiusKey = "max_radius_incidents"
	defaultMaxRadius     = 10000
	casheTTL             = 60 * time.Second
)

type Repository struct {
	redis *redis.Client
}

func NewIncidentCacheRepo(redis *infraredis.Client) *Repository {
	return &Repository{
		redis: redis.GetClient(),
	}
}
