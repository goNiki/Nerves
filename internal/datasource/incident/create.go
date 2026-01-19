package incident

import (
	"context"

	"github.com/goNiki/Nerves/internal/models"
	repoConverter "github.com/goNiki/Nerves/internal/repository/converter"
)

// записывает инцидент в Postgres и сразу добавляет его в кеш Redis
func (d *DataSource) CreateIncident(ctx context.Context, incident *models.Incident) (*models.Incident, error) {
	// добавление инцидента в PostgreSQL
	created, err := d.incidentRepo.Create(ctx, repoConverter.IncidentToRow(incident))
	if err != nil {
		return nil, err
	}

	// активные события сразу добавляем в кэш
	if created.Active {

		if err := d.cacheRepo.SetIncident(ctx, created); err != nil {
			d.logger.Error("failed to update cache after create", "id", created.ID, "error", err)
		}

		// обновляем максимальный радиус
		go func() {
			if err := d.updateMaxRadius(context.Background(), created.RadiusMeters); err != nil {
				d.logger.Warn("failed to update max radius in background", "error", err)
			}
		}()
	}

	d.logger.Info("incident created and cached", "id", created.ID)
	return created, nil
}

// обновляет максимальный радиус инцидента в кэше
func (d *DataSource) updateMaxRadius(ctx context.Context, newRadius int) error {

	currentMaxRadius, err := d.cacheRepo.GetMaxRadiusIncident(ctx)
	if err != nil {
		d.logger.Warn("failed to get max radius from cache", "error", err)
	}

	if currentMaxRadius == 0 {
		d.logger.Debug("max radius not in cache, fetching from DB")
		currentMaxRadius, err = d.incidentRepo.GetMaxRadiusMeters(ctx)
		if err != nil {
			return err
		}
	}

	if newRadius > currentMaxRadius {
		d.logger.Info("updating max radius", "old", currentMaxRadius, "new", newRadius)
		if err := d.cacheRepo.UpdateMaxRadiusIncident(ctx, float64(newRadius)); err != nil {
			return err
		}
	} else {
		d.logger.Debug("max radius not changed", "current", currentMaxRadius, "new", newRadius)
	}

	return nil
}
