package incident

import (
	"context"

	"github.com/google/uuid"
)

func (d *DataSource) DeactivateIncident(ctx context.Context, id uuid.UUID) error {
	// деактивация инцидента в БД(postgres)
	if err := d.incidentRepo.Deactivate(ctx, id); err != nil {
		return err
	}

	if err := d.cacheRepo.DeactivateIncident(ctx, id.String()); err != nil {
		d.logger.Warn("failed to delete deactive incident ID: %v: error: %w", id.String(), err)
	}
	return nil
}
