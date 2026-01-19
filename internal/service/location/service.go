package location

import (
	"context"
	"log/slog"

	"github.com/goNiki/Nerves/internal/datasource/location"
	"github.com/goNiki/Nerves/internal/models"
	"github.com/goNiki/Nerves/internal/repository"
)

type Service struct {
	locationDataSource *location.DataSource
	locationCheckRepo  repository.LocationCheckRepository
	webhookService     WebhookEnqueuer
	logger             *slog.Logger
}

type WebhookEnqueuer interface {
	Enqueue(ctx context.Context, task *models.WebhookTask) error
}

func New(
	locationDataSource *location.DataSource,
	locationCheckRepo repository.LocationCheckRepository,
	webhookService WebhookEnqueuer,
	logger *slog.Logger,
) *Service {
	return &Service{
		locationDataSource: locationDataSource,
		locationCheckRepo:  locationCheckRepo,
		webhookService:     webhookService,
		logger:             logger,
	}
}
