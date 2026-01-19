package incident

import (
	"context"

	"github.com/goNiki/Nerves/internal/models"
	repoConverter "github.com/goNiki/Nerves/internal/repository/converter"
)

// UpdateIncident используется для частичного обновления полей в БД
func (d *DataSource) UpdateIncident(ctx context.Context, update *models.UpdateIncident) (*models.Incident, error) {
	// получаем инцидент из БД
	existing, err := d.incidentRepo.GetById(ctx, update.ID)
	if err != nil {
		return nil, err
	}

	// заменяем nil поля данными из БД
	merged := d.mergeUpdates(existing, update)

	// сохрянаем в БД
	updated, err := d.incidentRepo.Update(ctx, repoConverter.IncidentToRow(merged))
	if err != nil {
		return nil, err
	}

	// обновление кеша
	if updated.Active {

		if err := d.cacheRepo.SetIncident(ctx, updated); err != nil {
			d.logger.Error("failed to update cache after patch ", "ID", updated.ID, "error", err)
		}

		// если поле радиуса менялось, то меняем кеш максимального инцидента
		if update.RadiusMeters != nil {
			go func() {
				if err := d.updateMaxRadius(context.Background(), updated.RadiusMeters); err != nil {
					d.logger.Warn("failed to update max radius in background", "error", err)
				}
			}()
		}
	} else {
		// если инцидент дектирован, то удаляем из БД
		if err := d.cacheRepo.DeactivateIncident(ctx, updated.ID.String()); err != nil {
			d.logger.Error("failed to delete cache incident", "id", updated.ID, "error", err)
		}
	}

	d.logger.Info("incident updated and cache synced", "id", updated.ID, "active", updated.Active)
	return updated, nil
}

// ReplaceIncident полное обновление всех полей инцидента
func (d *DataSource) ReplaceIncident(ctx context.Context, incident *models.Incident) (*models.Incident, error) {
	// 1. Обновляем в PostgreSQL
	updated, err := d.incidentRepo.Update(ctx, repoConverter.IncidentToRow(incident))
	if err != nil {
		return nil, err
	}

	if updated.Active {
		//обновляе кеш данного инцидента
		if err := d.cacheRepo.SetIncident(ctx, updated); err != nil {
			d.logger.Warn("failed to update cashe after replact", "id", updated.ID, "error", err)
		}
		// обновляем максимальный радиус инцидента
		go func() {
			if err := d.updateMaxRadius(context.Background(), incident.RadiusMeters); err != nil {
				d.logger.Warn("failed to update max radius incident", "error", err)
			}
		}()

	} else {
		// если обновленный инцидент не активный, то просто удаляем его.
		if err := d.cacheRepo.DeactivateIncident(ctx, updated.ID.String()); err != nil {
			d.logger.Warn("failed to delete deactive incident ", "ID:", updated.ID, "error", err)
		}
	}

	d.logger.Info("incident replaced and cache updated", "id", updated.ID, "active", updated.Active)
	return updated, nil
}

func (d *DataSource) mergeUpdates(existing *models.Incident, update *models.UpdateIncident) *models.Incident {
	merged := &models.Incident{
		ID:            update.ID,
		CreateAt:      existing.CreateAt,
		UpdatedAt:     existing.UpdatedAt,
		DeactivatedAt: existing.DeactivatedAt,
	}

	if update.Title != nil {
		merged.Title = *update.Title
	} else {
		merged.Title = existing.Title
	}

	if update.Description != nil {
		merged.Description = *update.Description
	} else {
		merged.Description = existing.Description
	}

	if update.Latitude != nil {
		merged.Latitude = *update.Latitude
	} else {
		merged.Latitude = existing.Latitude
	}

	if update.Longitude != nil {
		merged.Longitude = *update.Longitude
	} else {
		merged.Longitude = existing.Longitude
	}

	if update.RadiusMeters != nil {
		merged.RadiusMeters = *update.RadiusMeters
	} else {
		merged.RadiusMeters = existing.RadiusMeters
	}

	if update.Severity != nil {
		merged.Severity = *update.Severity
	} else {
		merged.Severity = existing.Severity
	}

	if update.Active != nil {
		merged.Active = *update.Active
	} else {
		merged.Active = existing.Active
	}

	return merged
}
