package cacheincident

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/goNiki/Nerves/internal/models"
	errorapp "github.com/goNiki/Nerves/internal/models/errorApp"
	"github.com/redis/go-redis/v9"
)

// Создание/обновление/деактивация одного инцидента
// если инцидент не создан, он будет создан, если есть, то будет изменен
func (r *Repository) SetIncident(ctx context.Context, incident *models.Incident) error {
	const op = "repository.cacheincident.setincident"

	data, err := json.Marshal(incident)
	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, errorapp.ErrMarshallCashIncident, err)
	}

	pipe := r.redis.Pipeline()

	pipe.HSet(ctx, activeIncidentsHash, incident.ID.String(), data)
	pipe.GeoAdd(ctx, incidentsGeo, &redis.GeoLocation{
		Name:      incident.ID.String(),
		Latitude:  incident.Latitude,
		Longitude: incident.Longitude,
	})

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("%s: %w: %v", op, errorapp.ErrSetToCache, err)
	}

	return nil
}

// загружает все активные инциденты в кеш
// использует атомарную замену через ReplaceAllActiveIncidents
func (r *Repository) SetAllActiveIncidents(ctx context.Context, incidents []*models.Incident) error {
	// используем ReplaceAllActiveIncidents для полной атомарной замены
	// это гарантирует что старые неактивные инциденты будут удалены
	return r.ReplaceAllActiveIncidents(ctx, incidents)
}

// полная замена хеша в двух таблицах на новый
// создаются новые, временные таблицы в redis, чтобы раньше времени не удалять хеш
// после загрузки всех инцидентов во временные таблицы, то создаем пайплайн в котором одновременно
// удаляем старый хеш и заменяем действующем.
func (r *Repository) ReplaceAllActiveIncidents(ctx context.Context, incidents []*models.Incident) error {
	const op = "repository.cacheincident.replaceallactiveincidents"

	tempHaskKey := activeIncidentsHash + ":temp"
	tempGeoKey := incidentsGeo + ":temp"

	pipe := r.redis.Pipeline()

	for _, incident := range incidents {
		if !incident.Active {
			continue
		}

		data, err := json.Marshal(incident)
		if err != nil {
			return fmt.Errorf("%s: %w: %w", op, errorapp.ErrMarshallCashIncident, err)
		}

		pipe.HSet(ctx, tempHaskKey, incident.ID.String(), data)

		pipe.GeoAdd(ctx, tempGeoKey, &redis.GeoLocation{
			Name:      incident.ID.String(),
			Latitude:  incident.Latitude,
			Longitude: incident.Longitude,
		})

	}

	if _, err := pipe.Exec(ctx); err != nil {
		r.redis.Del(ctx, tempHaskKey, tempGeoKey)
		return fmt.Errorf("%s: %w: %v", op, errorapp.ErrSetToCache, err)
	}

	pipe = r.redis.Pipeline()
	pipe.Del(ctx, activeIncidentsHash, incidentsGeo)
	pipe.Rename(ctx, tempHaskKey, activeIncidentsHash)
	pipe.Rename(ctx, tempGeoKey, incidentsGeo)

	if _, err := pipe.Exec(ctx); err != nil {
		r.redis.Del(ctx, tempHaskKey, tempGeoKey)
		return fmt.Errorf("%s: failed to replace hash: %v", op, err)
	}

	return nil
}

// очещает кеш всех активных инцидентов
func (r *Repository) InvalidateAllActiveIncidents(ctx context.Context) error {
	const op = "repository.cacheincident.invalidateallactiveincidents"

	if err := r.redis.Del(ctx, activeIncidentsHash, incidentsGeo).Err(); err != nil && err != redis.Nil {
		return fmt.Errorf("%s: %w: %v", op, errorapp.ErrDeleteCashIncident, err)
	}

	return nil
}

func (r *Repository) UpdateMaxRadiusIncident(ctx context.Context, maxRadius float64) error {
	const op = "repository.cacheincident.UpdateMaxRadiusIncident"

	err := r.redis.Set(ctx, maxIncidentRadiusKey, maxRadius, 10*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorapp.ErrSetToCache, err)
	}

	return nil
}
