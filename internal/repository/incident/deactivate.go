package incident

import (
	"context"
	"fmt"

	errorapp "github.com/goNiki/Nerves/internal/models/errorApp"
	"github.com/google/uuid"
)

func (r *Repository) Deactivate(ctx context.Context, id uuid.UUID) error {
	const op = "repository.incident.deactivate"

	query := `
			UPDATE incidents 
			SET active = false,
				updated_at = NOW(),
				deactivated_at = NOW()
			WHERE id = $1
	`

	result, err := r.db.DB.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %v: %w", op, errorapp.ErrDeactivateIncident, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorapp.ErrIncidentNotFound)
	}

	return nil

}
