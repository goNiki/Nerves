package cacheincident

import (
	"context"
	"fmt"

	errorapp "github.com/goNiki/Nerves/internal/models/errorApp"
)

func (r *Repository) DeactivateIncident(ctx context.Context, id string) error {
	const op = "repository.cacheincident.deactivateincident"

	pipe := r.redis.Pipeline()

	pipe.HDel(ctx, activeIncidentsHash, id)
	pipe.ZRem(ctx, incidentsGeo, id)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("%s: %w: %v", op, errorapp.ErrDeleteCashIncident, err)
	}
	return nil
}
