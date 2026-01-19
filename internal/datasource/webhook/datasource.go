package webhook

import (
	"log/slog"

	"github.com/goNiki/Nerves/internal/repository"
)

type DataSource struct {
	queueRepo repository.QueueRepository
	logger    *slog.Logger
}

func New(queueRepo repository.QueueRepository, logger *slog.Logger) *DataSource {
	return &DataSource{
		queueRepo: queueRepo,
		logger:    logger,
	}
}
