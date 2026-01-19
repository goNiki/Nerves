package cacheincident

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goNiki/Nerves/internal/models"
	errorapp "github.com/goNiki/Nerves/internal/models/errorApp"
	"github.com/goNiki/Nerves/internal/utils/geo"
	"github.com/redis/go-redis/v9"
)

// GetIncidentsWithDistanceInfo возвращает инциденты с информацией о расстояниях
func (r *Repository) GetIncidentsWithDistanceInfo(ctx context.Context, latitude, longitude float64, userRadiusMeters int) ([]*models.IncidentWithDistanceInfo, error) {
	const op = "repository.cacheincident.getincidentswithdistanceinfo"

	// Получаем максимальный радиус инцидента из Redis
	maxIncidentRadius := defaultMaxRadius
	if val, err := r.redis.Get(ctx, maxIncidentRadiusKey).Int(); err == nil {
		maxIncidentRadius = val
	}

	// Ищем инциденты в расширенном радиусе: userRadius + maxIncidentRadius
	searchRadius := float64(userRadiusMeters + maxIncidentRadius)

	result, err := r.redis.GeoRadius(ctx, incidentsGeo, longitude, latitude, &redis.GeoRadiusQuery{
		Radius:    searchRadius,
		Unit:      "m",
		WithCoord: true,
		WithDist:  true,
		Count:     1000,
		Sort:      "ASC",
	}).Result()

	if redis.Nil == err || len(result) == 0 {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorapp.ErrGetCacheAllActiveIncidents, err)
	}

	// Получаем ID инцидентов
	ids := make([]string, len(result))
	for i, v := range result {
		ids[i] = v.Name
	}

	// Получаем полные данные инцидентов из хеша
	datalist, err := r.redis.HMGet(ctx, activeIncidentsHash, ids...).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorapp.ErrGetCacheAllActiveIncidents, err)
	}

	incidents := make([]*models.IncidentWithDistanceInfo, 0, len(datalist))

	for _, data := range datalist {
		if data == nil {
			continue
		}

		var incident models.Incident
		if err := json.Unmarshal([]byte(data.(string)), &incident); err != nil {
			continue
		}

		// Вычисляем расстояния
		distanceToCenter := geo.CalculateDistance(latitude, longitude, incident.Latitude, incident.Longitude)
		distanceToIncidentBorder := distanceToCenter - float64(incident.RadiusMeters)
		isUserInside := distanceToCenter <= float64(incident.RadiusMeters)
		isInUserRange := distanceToCenter <= float64(userRadiusMeters)

		// Проверяем, попадает ли инцидент в зону интереса:
		// distance <= userRadius + incidentRadius
		if distanceToCenter <= float64(userRadiusMeters+incident.RadiusMeters) {
			incidents = append(incidents, &models.IncidentWithDistanceInfo{
				Incident:                 &incident,
				DistanceToCenter:         distanceToCenter,
				DistanceToIncidentBorder: distanceToIncidentBorder,
				IsUserInside:             isUserInside,
				IsInUserRange:            isInUserRange,
			})
		}
	}

	return incidents, nil
}
