package location

import (
	"context"
	"fmt"
	"time"

	"github.com/goNiki/Nerves/internal/models"
	errorapp "github.com/goNiki/Nerves/internal/models/errorApp"
	repoConverter "github.com/goNiki/Nerves/internal/repository/converter"
	repoModels "github.com/goNiki/Nerves/internal/repository/models"
)

func (r *Repository) GetUniqueUsersPerIncident(ctx context.Context, since time.Time) (*models.IncidentStatsList, error) {
	const op = "repository.location.getuniqueusersperincident"

	query := `
		SELECT
			i.id as incident_id,
			i.title,
			COUNT(DISTINCT lc.user_id) as user_count
		FROM incidents i
		LEFT JOIN location_checks lc ON i.id = ANY(lc.incident_ids)
		AND lc.checked_at >= $1
		WHERE i.active = true
		GROUP BY i.id, i.title
		ORDER BY user_count DESC
	`

	rows, err := r.db.DB.Query(ctx, query, since)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorapp.ErrGetIncidentStats, err)
	}

	defer rows.Close()

	var stats []repoModels.IncidentUserCountRow

	for rows.Next() {
		var row repoModels.IncidentUserCountRow

		if err := rows.Scan(
			&row.IncidentID,
			&row.Title,
			&row.UserCount,
		); err != nil {
			return nil, fmt.Errorf("%s: %w: %w", op, errorapp.ErrScanIcidentStats, err)
		}

		stats = append(stats, row)
	}

	return repoConverter.IncidentStatsLisToModels(stats), nil

}
