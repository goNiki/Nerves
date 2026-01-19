package cacheincident

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goNiki/Nerves/internal/models"
	errorapp "github.com/goNiki/Nerves/internal/models/errorApp"
	"github.com/goNiki/Nerves/internal/utils/geo"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// получает все активные инциденты из кеша
func (r *Repository) GetActiveIncidents(ctx context.Context) ([]*models.Incident, error) {
	const op = "repository.cacheincident.getactiveincidents"

	date, err := r.redis.HGetAll(ctx, activeIncidentsHash).Result()

	if redis.Nil == err || len(date) == 0 {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorapp.ErrGetCacheAllActiveIncidents, err)
	}

	incidents := make([]*models.Incident, 0, len(date))

	for _, v := range date {
		var incident models.Incident

		if err := json.Unmarshal([]byte(v), &incident); err != nil {
			continue
		}

		incidents = append(incidents, &incident)
	}

	return incidents, nil
}

// получает один инцидент из кеша по ID
func (r *Repository) GetIncident(ctx context.Context, id uuid.UUID) (*models.Incident, error) {
	const op = "repository.cacheincident.getincident"

	date, err := r.redis.HGet(ctx, activeIncidentsHash, id.String()).Result()
	if redis.Nil == err || len(date) == 0 {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorapp.ErrGetCacheActiveIncident, err)
	}

	var incident models.Incident

	if err := json.Unmarshal([]byte(date), &incident); err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorapp.ErrUnmarshallCashIncident, err)
	}

	return &incident, nil
}

func (r *Repository) GeoSearchIncidentIDs(ctx context.Context, latitude, longitude, searcheRadius float64) ([]string, error) {
	const op = "repository.cacheincident.geosearchincidentids"

	result, err := r.redis.GeoRadius(ctx, incidentsGeo, longitude, latitude, &redis.GeoRadiusQuery{
		Radius:    searcheRadius,
		Unit:      "m",
		WithCoord: false,
		WithDist:  false,
		Count:     10000,
		Sort:      "ASC",
	}).Result()

	if redis.Nil == err || len(result) == 0 {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorapp.ErrGeoSearchIncidentIDs, err)
	}

	ids := make([]string, len(result))
	for i, v := range result {
		ids[i] = v.Name
	}
	return ids, nil
}

func (r *Repository) GetIncidentsByIDs(ctx context.Context, ids []string) ([]*models.Incident, error) {
	const op = "repository.cacheincident.getincidentsbyids"

	datalist, err := r.redis.HMGet(ctx, activeIncidentsHash, ids...).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorapp.ErrGetIncidentsByIDs, err)
	}

	incidents := make([]*models.Incident, 0, len(datalist))

	for _, date := range datalist {
		if date == nil {
			continue
		}

		var incident models.Incident

		if err := json.Unmarshal([]byte(date.(string)), &incident); err != nil {
			continue
		}

		incidents = append(incidents, &incident)
	}
	return incidents, nil
}

// скорее всего нужно будет удалить.
func (r *Repository) GetIncidentsWithRadius(ctx context.Context, latitude, longitude float64, userRadiusMeters int) ([]*models.Incident, error) {
	const op = "repository.cacheincident.getincidentswithradius"

	// Получаем максимальный радиус инцидента из Redis
	maxIncidentRadius := defaultMaxRadius
	if val, err := r.redis.Get(ctx, maxIncidentRadiusKey).Int(); err == nil {
		maxIncidentRadius = val
	}

	// Ищем инциденты в расширенном радиусе: userRadius + maxIncidentRadius
	// Это гарантирует, что мы найдем все инциденты, в радиус которых может попасть пользователь
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

	incidents := make([]*models.Incident, 0, len(datalist))

	for _, data := range datalist {
		if data == nil {
			continue
		}

		var incident models.Incident
		if err := json.Unmarshal([]byte(data.(string)), &incident); err != nil {
			continue
		}

		// Вычисляем точное расстояние до центра инцидента
		distanceToCenter := geo.CalculateDistance(latitude, longitude, incident.Latitude, incident.Longitude)

		// Проверяем, попадает ли инцидент в зону интереса:
		// Пользователь должен быть либо в радиусе инцидента, либо инцидент в радиусе пользователя
		if distanceToCenter <= float64(userRadiusMeters+incident.RadiusMeters) {
			incidents = append(incidents, &incident)
		}
	}

	return incidents, nil
}

func (r *Repository) GetMaxRadiusIncident(ctx context.Context) (int, error) {
	const op = "repository.cacheincinedt.getmaxradiusincidents"

	result, err := r.redis.Get(ctx, maxIncidentRadiusKey).Int()
	if redis.Nil == err {
		return 0, nil
	}

	if err != nil {
		return 0, fmt.Errorf("%s: %w: %w", op, errorapp.ErrMaxRadiusNotFound)
	}

	return result, nil

}
