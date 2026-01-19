package incident

import (
	"log/slog"

	"github.com/goNiki/Nerves/internal/repository"
)

// DataSource - прослойка между Service и Repository
// Реализует паттерн Cache-Aside для работы с инцидентами
type DataSource struct {
	incidentRepo repository.IncidentRepository
	cacheRepo    repository.CacheRepository
	logger       *slog.Logger
}

func NewDataSource(
	incidentRepo repository.IncidentRepository,
	cacheRepo repository.CacheRepository,
	logger *slog.Logger,
) *DataSource {
	return &DataSource{
		incidentRepo: incidentRepo,
		cacheRepo:    cacheRepo,
		logger:       logger,
	}
}
