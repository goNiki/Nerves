package cache

import (
	"context"
	"log/slog"
	"time"

	"github.com/goNiki/Nerves/internal/datasource/incident"
)

type Worker struct {
	incidentDataSource *incident.DataSource
	refreshInterval    time.Duration
	logger             *slog.Logger
}

func New(
	incidentDataSource *incident.DataSource,
	refreshInterval time.Duration,
	logger *slog.Logger,
) *Worker {
	return &Worker{
		incidentDataSource: incidentDataSource,
		refreshInterval:    refreshInterval,
		logger:             logger,
	}
}

func (w *Worker) Run(ctx context.Context) {
	w.logger.Info("cache refresh worker started", "interval", w.refreshInterval)

	ticker := time.NewTicker(w.refreshInterval)
	defer ticker.Stop()

	// первое обновление сразу при старте
	w.refreshCache(ctx)

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("cache refresh worker stopped")
			return
		case <-ticker.C:
			w.refreshCache(ctx)
		}
	}
}

func (w *Worker) refreshCache(ctx context.Context) {
	start := time.Now()

	w.logger.Info("starting cache refresh")

	// принудительно обновляем кэш активных инцидентов
	incidents, err := w.incidentDataSource.RefreshActiveIncidentsCache(ctx)
	if err != nil {
		w.logger.Error("failed to refresh active incidents cache", "error", err)
		return
	}

	duration := time.Since(start)

	w.logger.Info("cache refresh completed",
		"incidents_count", len(incidents),
		"duration", duration,
	)
}
