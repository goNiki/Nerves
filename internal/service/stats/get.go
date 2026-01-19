package stats

import (
	"context"
	"time"

	"github.com/goNiki/Nerves/internal/models"
)

func (s *Service) GetIncidentStats(ctx context.Context) (*models.IncidentStatsList, error) {
	since := time.Now().Add(-time.Duration(s.timeWindowMinutes) * time.Minute)

	stats, err := s.locationCheckRepo.GetUniqueUsersPerIncident(ctx, since)
	if err != nil {
		s.logger.Error("failed to get incident stats", "error", err)
		return nil, err
	}

	return stats, nil
}
