package location

import (
	"context"

	"github.com/goNiki/Nerves/internal/models"
	"github.com/goNiki/Nerves/internal/utils/geo"
)

func (d *DataSource) Check(ctx context.Context, latitude, longitude float64, userRadiusMeters int) ([]*models.IncidentWithDistanceInfo, error) {
	// получаем максимальный радиус инцидента с кеша
	maxRadius := d.getMaxRadiusIncident(ctx)

	// получаем радиус поиска, путем складывания радиус инцидента и радиус юзера
	// необходимо чтобы понимать, когда начинается радиус действие инцидента
	searchRadius := float64(maxRadius + userRadiusMeters)

	ids, err := d.cacheRepo.GeoSearchIncidentIDs(ctx, latitude, longitude, searchRadius)
	if err != nil {
		return nil, err
	}

	var incidents []*models.Incident

	//если id инцидентов нашлись, то мы получаем всю информацию об инцидентах
	if ids != nil && len(ids) > 0 {
		d.logger.Debug("geo index hit", "count", len(ids))

		incidents, err = d.cacheRepo.GetIncidentsByIDs(ctx, ids)
		if err != nil {
			return nil, err
		}
	}

	// если через geo индекс ничего не нашлось - проверяем пытаемся найти инциденты вручную
	if ids == nil || len(ids) == 0 {
		d.logger.Debug("geo index empty, trying to get all active incidents")

		incidents, err = d.getIncidentsInRadius(ctx, latitude, longitude, searchRadius)
		if err != nil {
			d.logger.Warn("failed to get incident", "error", err)
			return nil, err
		}

	}
	if incidents == nil || len(incidents) == 0 {
		d.logger.Debug("no incidents found in area")
		return nil, nil
	}

	result := d.buildIncidentsWithDistance(incidents, latitude, longitude, userRadiusMeters)

	d.logger.Info("location check completed",
		"lat", latitude,
		"lon", longitude,
		"incidents_found", len(result))

	return result, nil
}

// getMaxRadiusIncident - позволяет получить max radius incident из кэша redis или из postgres
func (d *DataSource) getMaxRadiusIncident(ctx context.Context) int {
	maxRadius, err := d.cacheRepo.GetMaxRadiusIncident(ctx)
	if err != nil || maxRadius == 0 {
		d.logger.Debug("cache miss maxradius")
		// cache miss - получаем данные из postgres
		maxRadius, err = d.incidentRepo.GetMaxRadiusMeters(ctx)
		if err != nil {
			d.logger.Warn("failed to get max radius", "error", err)
			// если не удалось получить, то используем дефолтный
			maxRadius = models.DefaultMaxRadius
		} else {
			go func() {
				if err := d.cacheRepo.UpdateMaxRadiusIncident(ctx, float64(maxRadius)); err != nil {
					d.logger.Warn("failed to update max radius cache", "error", err)
				}
				d.logger.Info("updated max radius cache", "maxradius", maxRadius)
			}()
		}
	}
	return maxRadius
}

func (d *DataSource) getIncidentsInRadius(ctx context.Context, latitude, longitede, searcheRadius float64) ([]*models.Incident, error) {
	//достаем все активные инциденты из кэша
	incidents, err := d.cacheRepo.GetActiveIncidents(ctx)
	if err != nil {
		d.logger.Warn("failed to get active incidents from cache", "error", err)
	}

	//если кэш пустой или произошел cache miss,
	//смотрим в Postgres, точно ли нет активных инцидентов
	if incidents == nil || len(incidents) == 0 {
		d.logger.Debug("cache miss: active incindets")

		//достаем все активные инциденты из БД
		incidents, err = d.incidentRepo.GetAllActiveWithoutLimit(ctx)
		if err != nil {
			return nil, err
		}

		if incidents != nil && len(incidents) > 0 {
			// если активные инциденты есть в базе, то начинам обновлять кэш ассинхронно
			go func() {
				if err := d.cacheRepo.SetAllActiveIncidents(context.Background(), incidents); err != nil {
					d.logger.Warn("failed to update cache with active incidents", "error", err)
				}
			}()
		}
	}
	// если в базе есть активные инциденты, то необходимо их отфильтровать по расстоянию
	if incidents != nil && len(incidents) > 0 {
		incidents = d.filterByDistance(incidents, latitude, longitede, searcheRadius)
		return incidents, nil
	}
	return nil, nil
}

func (d *DataSource) filterByDistance(incidents []*models.Incident, latitude, longitude, maxDistance float64) []*models.Incident {
	//ставлю пока 5% от общего количества активных инцидентов
	result := make([]*models.Incident, len(incidents)*0, 05)

	for _, inc := range incidents {
		distance := geo.CalculateDistance(latitude, longitude, inc.Latitude, inc.Longitude)
		if distance <= maxDistance {
			result = append(result, inc)
		}
	}
	//ставим метку дебага, чтобы потом проставить корректную длину слайса,
	// чтобы постоянно не переалацировать его
	d.logger.Debug("filtered by distance", "before", len(incidents), "after", len(result))

	return result
}

func (d *DataSource) buildIncidentsWithDistance(incidents []*models.Incident, latitude, longitude float64, userRadiusMeters int) []*models.IncidentWithDistanceInfo {
	result := make([]*models.IncidentWithDistanceInfo, 0, len(incidents))

	for _, inc := range incidents {
		//находим расстояние от человека до центра инцидента
		distanceToCenter := geo.CalculateDistance(latitude, longitude, inc.Latitude, inc.Longitude)

		if distanceToCenter <= float64(userRadiusMeters+inc.RadiusMeters) {
			result = append(result, &models.IncidentWithDistanceInfo{
				Incident:                 inc,
				DistanceToCenter:         distanceToCenter,
				DistanceToIncidentBorder: distanceToCenter - float64(inc.RadiusMeters),
				IsUserInside:             distanceToCenter <= float64(inc.RadiusMeters),
				IsInUserRange:            distanceToCenter <= float64(userRadiusMeters),
			})
		}
	}

	return result
}
