package location

import (
	"log/slog"

	"github.com/goNiki/Nerves/internal/repository"
)

type DataSource struct {
	cacheRepo    repository.CacheRepository
	incidentRepo repository.IncidentRepository
	locationRepo repository.LocationCheckRepository
	logger       *slog.Logger
}

func New(
	cacheRepo repository.CacheRepository,
	incidentRepo repository.IncidentRepository,
	logger *slog.Logger,
) *DataSource {
	return &DataSource{
		cacheRepo:    cacheRepo,
		incidentRepo: incidentRepo,
		logger:       logger,
	}
}
