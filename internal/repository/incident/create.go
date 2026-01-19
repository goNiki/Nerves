package incident

import (
	"context"
	"fmt"

	"github.com/goNiki/Nerves/internal/models"
	errorapp "github.com/goNiki/Nerves/internal/models/errorApp"
	repoConverter "github.com/goNiki/Nerves/internal/repository/converter"
	repoModels "github.com/goNiki/Nerves/internal/repository/models"
)

func (r *Repository) Create(ctx context.Context, incident *repoModels.IncidentRow) (*models.Incident, error) {
	const op = "repository.incident.create"

	query := `
		INSERT INTO incidents(title, description, latitude, longitude, radius_meters, severity, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`
	err := r.db.DB.QueryRow(ctx, query,
		incident.Title,
		incident.Description,
		incident.Latitude,
		incident.Longitude,
		incident.RadiusMeters,
		incident.Severity,
		incident.Active).Scan(&incident.ID, &incident.CreateAt)

	if err != nil {
		return nil, fmt.Errorf("%s: %v :%w", op, errorapp.ErrCreateIncident, err)
	}

	return repoConverter.RowToIncident(incident), nil

}
