package incident

import (
	"log/slog"

	"github.com/goNiki/Nerves/internal/datasource/incident"
)

type Service struct {
	dataSource *incident.DataSource
	logger     *slog.Logger
}

func NewIncidentService(dataSource *incident.DataSource, logger *slog.Logger) *Service {
	return &Service{
		dataSource: dataSource,
		logger:     logger,
	}
}
