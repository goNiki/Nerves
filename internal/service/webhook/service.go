package webhook

import (
	"log/slog"

	"github.com/goNiki/Nerves/internal/datasource/webhook"
)

type Service struct {
	dataSource *webhook.DataSource
	logger     *slog.Logger
}

func New(dataSource *webhook.DataSource, logger *slog.Logger) *Service {
	return &Service{
		dataSource: dataSource,
		logger:     logger,
	}
}
