package incident

import (
	"context"

	"github.com/google/uuid"
)

func (s *Service) Deactivate(ctx context.Context, id uuid.UUID) error {

	if err := s.dataSource.DeactivateIncident(ctx, id); err != nil {
		s.logger.Error("failed to deactivate incident", "id", id, "error", err)
		return err
	}

	s.logger.Info("incident deactivated", "id", id)
	return nil
}
