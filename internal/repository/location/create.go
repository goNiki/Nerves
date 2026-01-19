package location

import (
	"context"
	"fmt"

	"github.com/goNiki/Nerves/internal/models"
	errorapp "github.com/goNiki/Nerves/internal/models/errorApp"
	repoConverter "github.com/goNiki/Nerves/internal/repository/converter"
	repoModels "github.com/goNiki/Nerves/internal/repository/models"
)

func (r *Repository) Create(ctx context.Context, check *repoModels.LocationCheckRow) (*models.LocationCheck, error) {
	const op = "repository.incident.create"

	query := `INSERT INTO location_checks (user_id, latitude, longitude, incidents_found, incident_ids)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, checked_at
	`

	err := r.db.DB.QueryRow(ctx, query,
		check.UserID,
		check.Latitude,
		check.Longitude,
		check.IncidentsFound,
		check.IncidentIDs).Scan(&check.ID, &check.CheckedAt)

	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorapp.ErrCreateLocationCheck, err)
	}

	return repoConverter.LocationCheckRowToModels(check), nil
}
