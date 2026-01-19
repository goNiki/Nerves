package service

import (
	"context"

	"github.com/goNiki/Nerves/internal/models"
	"github.com/google/uuid"
)

type IncidentService interface {
	Create(ctx context.Context, incident *models.Incident) (*models.Incident, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Incident, error)
	List(ctx context.Context, page, limit int) (*models.IncidentList, error)
	Replace(ctx context.Context, id uuid.UUID, incident *models.Incident) (*models.Incident, error)
	Update(ctx context.Context, id uuid.UUID, updates *models.UpdateIncident) (*models.Incident, error)
	Deactivate(ctx context.Context, id uuid.UUID) error
}

type LocationService interface {
	Check(ctx context.Context, req *models.LocationCheckRequest) (*models.LocationCheckResponce, error)
}

type StatsService interface {
	GetIncidentStats(ctx context.Context) (*models.IncidentStatsList, error)
}

type HealthService interface {
	Check(ctx context.Context) (*models.HealthStatus, error)
}

type WebhookService interface {
	Enqueue(ctx context.Context, task *models.WebhookTask) error
}
