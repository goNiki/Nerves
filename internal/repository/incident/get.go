package incident

import (
	"context"
	"fmt"

	"github.com/goNiki/Nerves/internal/models"
	errorapp "github.com/goNiki/Nerves/internal/models/errorApp"
	repoConverter "github.com/goNiki/Nerves/internal/repository/converter"
	repoModels "github.com/goNiki/Nerves/internal/repository/models"
	"github.com/google/uuid"
)

// добавить поле деактивации
func (r *Repository) GetById(ctx context.Context, id uuid.UUID) (*models.Incident, error) {
	const op = "repository.incident.getbyid"

	query := `
		SELECT id, title, description, latitude, longitude, radius_meters, severity, active, created_at, updated_at, deactivated_at
		FROM incidents
		WHERE id = $1
	`

	var result repoModels.IncidentRow

	err := r.db.DB.QueryRow(ctx, query, id).Scan(
		&result.ID,
		&result.Title,
		&result.Description,
		&result.Latitude,
		&result.Longitude,
		&result.RadiusMeters,
		&result.Severity,
		&result.Active,
		&result.CreateAt,
		&result.UpdatedAt,
		&result.DeactivatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %v :%w", op, errorapp.ErrGetIncident, err)
	}

	return repoConverter.RowToIncident(&result), nil

}

func (r *Repository) GetListIncident(ctx context.Context, offset, limit int) ([]*models.Incident, int, error) {
	const op = "repository.incident.getlistincident"

	var total int
	//уточнить, нужно прям все инциденты искать или только те, которые находятся в радиусе какой-то координаты
	countQuery := `SELECT COUNT(*) FROM incidents`
	if err := r.db.DB.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("%s: %v: %w", op, errorapp.ErrCountIncidents, err)
	}

	query := `
		SELECT id, title, description, latitude, longitude, radius_meters, severity, active, created_at, updated_at, deactivated_at
		FROM incidents
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.DB.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %v: %w", op, errorapp.ErrListIncidents, err)
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
			return nil, 0, fmt.Errorf("%s: %v: %w", op, errorapp.ErrScanIncident, err)
		}
		incidents = append(incidents, &row)
	}
	return repoConverter.RowsToIncidents(incidents), total, nil
}

func (r *Repository) GetAllActive(ctx context.Context, offset, limit int) ([]*models.Incident, int, error) {
	const op = "repository.incident.getallactive"

	var totalActive int
	countQuery := `SELECT COUNT(*) FROM incidents WHERE active = true`

	if err := r.db.DB.QueryRow(ctx, countQuery).Scan(&totalActive); err != nil {
		return nil, 0, fmt.Errorf("%s: %w: %w", op, errorapp.ErrCountIncidents, err)
	}

	query := `
		SELECT id, title, description, latitude, 
		longitude, radius_meters, severity, active, 
		created_at, updated_at, deactivated_at
		FROM incidents
		WHERE active = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.DB.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w: %w", op, errorapp.ErrListIncidents, err)
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
			return nil, 0, fmt.Errorf("%s: %w: %w", op, errorapp.ErrScanIncident, err)
		}

		incidents = append(incidents, &row)
	}

	return repoConverter.RowsToIncidents(incidents), totalActive, nil
}
