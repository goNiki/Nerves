package incident

import (
	"context"

	"github.com/goNiki/Nerves/internal/models"
)

// RefreshActiveIncidentsCache принудительно обновляет кэш активных инцидентов
// Используется cache refresh worker для периодического обновления
func (d *DataSource) RefreshActiveIncidentsCache(ctx context.Context) ([]*models.Incident, error) {
	d.logger.Debug("force refreshing active incidents cache")

	// получаем все активные инциденты из БД
	incidents, err := d.incidentRepo.GetAllActiveWithoutLimit(ctx)
	if err != nil {
		d.logger.Error("failed to get active incidents from db", "error", err)
		return nil, err
	}

	if len(incidents) == 0 {
		d.logger.Debug("no active incidents found")
		return incidents, nil
	}

	// обновляем кэш активных инцидентов (включая geo-индекс)
	if err := d.cacheRepo.SetAllActiveIncidents(ctx, incidents); err != nil {
		d.logger.Error("failed to update active incidents cache", "error", err)
		return nil, err
	}

	// обновляем максимальный радиус
	maxRadius := d.calculateMaxRadius(incidents)
	if err := d.cacheRepo.UpdateMaxRadiusIncident(ctx, float64(maxRadius)); err != nil {
		d.logger.Warn("failed to update max radius cache", "error", err)
	}

	d.logger.Info("cache refreshed successfully",
		"incidents_count", len(incidents),
		"max_radius", maxRadius,
	)

	return incidents, nil
}

// calculateMaxRadius находит максимальный радиус среди инцидентов
func (d *DataSource) calculateMaxRadius(incidents []*models.Incident) int {
	maxRadius := 0
	for _, inc := range incidents {
		if inc.RadiusMeters > maxRadius {
			maxRadius = inc.RadiusMeters
		}
	}
	return maxRadius
}
