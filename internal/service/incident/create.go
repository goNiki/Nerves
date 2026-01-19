package incident

import (
	"context"

	"github.com/goNiki/Nerves/internal/models"
)

func (s *Service) Create(ctx context.Context, incident *models.Incident) (*models.Incident, error) {

	created, err := s.dataSource.CreateIncident(ctx, incident)
	if err != nil {
		s.logger.Error("failed to create incident", "error", err)
		return nil, err
	}

	s.logger.Info("incident created", "id", created.ID)

	return created, nil
}
