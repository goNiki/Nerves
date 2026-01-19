package health

import (
	"log/slog"

	"github.com/goNiki/Nerves/internal/infrastructure/database"
	infraredis "github.com/goNiki/Nerves/internal/infrastructure/redis"
)

type Service struct {
	db     *database.Postgres
	redis  *infraredis.Client
	logger *slog.Logger
}

func New(db *database.Postgres, redis *infraredis.Client, logger *slog.Logger) *Service {
	return &Service{
		db:     db,
		redis:  redis,
		logger: logger,
	}
}
