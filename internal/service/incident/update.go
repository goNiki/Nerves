package incident

import (
	"context"

	"github.com/goNiki/Nerves/internal/models"
	"github.com/google/uuid"
)

func (s *Service) Replace(ctx context.Context, id uuid.UUID, incident *models.Incident) (*models.Incident, error) {

	incident.ID = id

	updated, err := s.dataSource.ReplaceIncident(ctx, incident)
	if err != nil {
		s.logger.Error("failed to replact incident", "id", id, "error", err)
		return nil, err
	}

	s.logger.Info("incident replaced", "id", id)
	return updated, nil
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, updates *models.UpdateIncident) (*models.Incident, error) {

	updates.ID = id

	updated, err := s.dataSource.UpdateIncident(ctx, updates)
	if err != nil {
		s.logger.Error("failed to update incident", "id", id, "error", err)
		return nil, err
	}

	s.logger.Info("incident updated", "id ", id)
	return updated, nil
}
