package repository

import (
	"context"
	"time"

	"github.com/goNiki/Nerves/internal/models"
	repoModels "github.com/goNiki/Nerves/internal/repository/models"
	"github.com/google/uuid"
)

// IncidentRepository интерфейс для работы с инцидентами в БД
type IncidentRepository interface {
	Create(ctx context.Context, incident *repoModels.IncidentRow) (*models.Incident, error)
	GetById(ctx context.Context, id uuid.UUID) (*models.Incident, error)
	GetListIncident(ctx context.Context, offset, limit int) ([]*models.Incident, int, error)
	GetAllActive(ctx context.Context, offset, limit int) ([]*models.Incident, int, error)
	GetAllActiveWithoutLimit(ctx context.Context) ([]*models.Incident, error)
	GetMaxRadiusMeters(ctx context.Context) (int, error)
	Update(ctx context.Context, row *repoModels.IncidentRow) (*models.Incident, error)
	Deactivate(ctx context.Context, id uuid.UUID) error
}

// LocationCheckRepository интерфейс для работы с проверками локаций
type LocationCheckRepository interface {
	Create(ctx context.Context, check *repoModels.LocationCheckRow) (*models.LocationCheck, error)
	GetUniqueUsersPerIncident(ctx context.Context, since time.Time) (*models.IncidentStatsList, error)
}

// CacheRepository интерфейс для работы с кэшем инцидентов в Redis
type CacheRepository interface {
	GetActiveIncidents(ctx context.Context) ([]*models.Incident, error)
	GetIncidentsByIDs(ctx context.Context, ids []string) ([]*models.Incident, error)
	GetIncident(ctx context.Context, id uuid.UUID) (*models.Incident, error)
	GetIncidentsWithRadius(ctx context.Context, latitude, longitude float64, userRadiusMeters int) ([]*models.Incident, error)
	GetIncidentsWithDistanceInfo(ctx context.Context, latitude, longitude float64, userRadiusMeters int) ([]*models.IncidentWithDistanceInfo, error)
	GeoSearchIncidentIDs(ctx context.Context, latitude, longitude, searcheRadius float64) ([]string, error)
	GetMaxRadiusIncident(ctx context.Context) (int, error)
	SetIncident(ctx context.Context, incident *models.Incident) error
	SetAllActiveIncidents(ctx context.Context, incidents []*models.Incident) error
	ReplaceAllActiveIncidents(ctx context.Context, incidents []*models.Incident) error
	InvalidateAllActiveIncidents(ctx context.Context) error
	UpdateMaxRadiusIncident(ctx context.Context, maxRadius float64) error
	DeactivateIncident(ctx context.Context, id string) error
}

// QueueRepository интерфейс для работы с очередью вебхуков в Redis
type QueueRepository interface {
	EnqueueWebhook(ctx context.Context, task *models.WebhookTask) error
	DequeueWebhook(ctx context.Context) (*models.WebhookTask, error)
}
