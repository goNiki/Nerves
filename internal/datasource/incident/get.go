package incident

import (
	"context"

	"github.com/goNiki/Nerves/internal/models"
	"github.com/google/uuid"
)

func (d *DataSource) GetByID(ctx context.Context, id uuid.UUID) (*models.Incident, error) {

	incident, err := d.cacheRepo.GetIncident(ctx, id)
	if err != nil {
		d.logger.Warn("cache error on GetByID", "id", id, "error", err)
	}

	//если произошел cashe hit
	if incident != nil {
		d.logger.Debug("cache hit", "id", id)
		return incident, nil
	}

	//если произошел cache miss, то идем в БД за инцидентом
	d.logger.Debug("cache miss", "id", id)
	incident, err = d.incidentRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := d.cacheRepo.SetIncident(context.Background(), incident); err != nil {
			d.logger.Warn("failed to update cache incident", "id", id, "error", err)
		}
	}()

	return incident, nil

}

func (d *DataSource) GetAllActive(ctx context.Context) ([]*models.Incident, error) {
	//получаем активные события их кэша Redis
	incidents, err := d.cacheRepo.GetActiveIncidents(ctx)
	if err != nil {
		d.logger.Warn("cache error on GetAllActive", "error", err)
	}

	//если произошел cache hit, то возвращаем значение
	if len(incidents) > 0 {
		d.logger.Debug("cache hit for active incidents", "count", len(incidents))
		return incidents, nil
	}

	// при cache miss идем в postgres
	d.logger.Debug("cache miss for active incidents")
	incidents, err = d.incidentRepo.GetAllActiveWithoutLimit(ctx)
	if err != nil {
		return nil, err
	}
	// записываем в кэш активные инциденты
	go func() {
		if err := d.cacheRepo.SetAllActiveIncidents(context.Background(), incidents); err != nil {
			d.logger.Warn("failed to update cache", "error", err)
		}
	}()

	return incidents, nil
}

func (d *DataSource) GetListIncident(ctx context.Context, offset, limit int) ([]*models.Incident, int, error) {
	return d.incidentRepo.GetListIncident(ctx, offset, limit)
}
