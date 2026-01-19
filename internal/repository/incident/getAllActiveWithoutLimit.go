package incident

import (
	"context"
	"fmt"

	"github.com/goNiki/Nerves/internal/models"
	errorapp "github.com/goNiki/Nerves/internal/models/errorApp"
	repoConverter "github.com/goNiki/Nerves/internal/repository/converter"
	repoModels "github.com/goNiki/Nerves/internal/repository/models"
)

// для получения всех активных инциденты без пагинации
func (r *Repository) GetAllActiveWithoutLimit(ctx context.Context) ([]*models.Incident, error) {
	const op = "repository.incident.getallactivewithoutlimit"

	query := `
		SELECT id, title, description, latitude, longitude, radius_meters, severity, active, 
		       created_at, updated_at, deactivated_at
		FROM incidents
		WHERE active = true
		ORDER BY created_at DESC
	`

	rows, err := r.db.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorapp.ErrListIncidents, err)
	}
	defer rows.Close()

	var incidents []*repoModels.IncidentRow
	for rows.Next() {
		var row repoModels.IncidentRow
		if err := rows.Scan(
			&row.ID,
			&row.Title,
			&row.Description,
			&row.Latitude,
			&row.Longitude,
			&row.RadiusMeters,
			&row.Severity,
			&row.Active,
			&row.CreateAt,
			&row.UpdatedAt,
			&row.DeactivatedAt,
		); err != nil {
			return nil, fmt.Errorf("%s: %w: %w", op, errorapp.ErrScanIncident, err)
		}

		incidents = append(incidents, &row)
	}

	return repoConverter.RowsToIncidents(incidents), nil
}
