package health

import (
	"context"

	"github.com/goNiki/Nerves/internal/models"
)

func (s *Service) Check(ctx context.Context) (*models.HealthStatus, error) {
	status := &models.HealthStatus{
		Status:     "healthy",
		Components: make(map[string]string),
	}

	// проверяем postgres
	if err := s.db.DB.Ping(ctx); err != nil {
		s.logger.Error("postgres health check failed", "error", err)
		status.Components["postgres"] = "down"
		status.Status = "unhealthy"
	} else {
		status.Components["postgres"] = "up"
	}

	// проверяем redis
	if err := s.redis.GetClient().Ping(ctx).Err(); err != nil {
		s.logger.Error("redis health check failed", "error", err)
		status.Components["redis"] = "down"
		status.Status = "unhealthy"
	} else {
		status.Components["redis"] = "up"
	}

	return status, nil
}
