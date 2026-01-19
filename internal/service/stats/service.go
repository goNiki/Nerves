package stats

import (
	"log/slog"

	"github.com/goNiki/Nerves/internal/config"
	"github.com/goNiki/Nerves/internal/repository"
)

type Service struct {
	locationCheckRepo repository.LocationCheckRepository
	timeWindowMinutes int
	logger            *slog.Logger
}

func New(
	locationCheckRepo repository.LocationCheckRepository, cfg config.Stats, logger *slog.Logger) *Service {
	return &Service{
		locationCheckRepo: locationCheckRepo,
		timeWindowMinutes: cfg.TimeWindiwMinutes(),
		logger:            logger,
	}
}
