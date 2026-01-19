package incident

import (
	"context"
	"fmt"

	"github.com/goNiki/Nerves/internal/models"
	errorapp "github.com/goNiki/Nerves/internal/models/errorApp"
	repoConverter "github.com/goNiki/Nerves/internal/repository/converter"
	repoModels "github.com/goNiki/Nerves/internal/repository/models"
)

// добавить проставление даты деактивации
func (r *Repository) Update(ctx context.Context, row *repoModels.IncidentRow) (*models.Incident, error) {
	const op = "repository.incident.update"

	query := `
			UPDATE Incidents
			SET title = $2,
				description = $3,
				latetude = $4,
				longetude = $5,
				radius_meters = $6,
				severity = $7,
				active = $8,
				updated_at = NOW()
			WHERE id = $1
			RETURNING updated_at
	`

	err := r.db.DB.QueryRow(ctx, query,
		row.ID,
		row.Title,
		row.Description,
		row.Latitude,
		row.Longitude,
		row.RadiusMeters,
		row.Severity,
		row.Active,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %v: %w", op, errorapp.ErrUpdateIncident, err)
	}

	return repoConverter.RowToIncident(row), nil

}
