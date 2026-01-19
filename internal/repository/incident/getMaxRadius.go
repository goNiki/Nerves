package incident

import (
	"context"
	"fmt"

	errorapp "github.com/goNiki/Nerves/internal/models/errorApp"
)

func (r *Repository) GetMaxRadiusMeters(ctx context.Context) (int, error) {
	const op = "repository.incident.getmaxradiusmeters"

	query := `
		SELECT COALESCE(MAX(radius_meters), 0) as max_radius
		FROM incidents
		WHERE active = true
	`

	var maxRadius int
	err := r.db.DB.QueryRow(ctx, query).Scan(&maxRadius)
	if err != nil {
		return 0, fmt.Errorf("%s: %w: %w", op, errorapp.ErrGetIncident, err)
	}

	return maxRadius, nil
}
