package location

import (
	"context"
	"time"

	"github.com/goNiki/Nerves/internal/models"
	repoModels "github.com/goNiki/Nerves/internal/repository/models"
	"github.com/google/uuid"
)

// основной метод системы - должен работать < 100ms

func (s *Service) Check(ctx context.Context, req *models.LocationCheckRequest) (*models.LocationCheckResponce, error) {
	const defaultUserSearchRadius = 5000 // 5 км

	// получаем инциденты с расстояниями через DataSource
	incidentsWithDistance, err := s.locationDataSource.Check(ctx, req.Latitude, req.Longitude, defaultUserSearchRadius)
	if err != nil {
		s.logger.Error("failed to check location", "error", err)
		return nil, err
	}

	// извлекаем только инциденты для ответа
	incidents := make([]*models.Incident, len(incidentsWithDistance))
	for i, iwd := range incidentsWithDistance {
		incidents[i] = iwd.Incident
	}

	response := &models.LocationCheckResponce{
		Incident:  incidents,
		CheckedAt: time.Now(),
	}

	// асинхронно сохраняем факт проверки и ставим задачу вебхука
	go func() {
		incidentIDs := make([]uuid.UUID, len(incidents))
		for i, inc := range incidents {
			incidentIDs[i] = inc.ID
		}

		locationCheck := &repoModels.LocationCheckRow{
			UserID:         req.UserID,
			Latitude:       req.Latitude,
			Longitude:      req.Longitude,
			IncidentsFound: len(incidents),
			IncidentIDs:    incidentIDs,
		}

		if _, err := s.locationCheckRepo.Create(context.Background(), locationCheck); err != nil {
			s.logger.Error("failed to save location check", "error", err)
		}

		if len(incidents) > 0 {
			task := models.NewWebhookTask(req.UserID, req.Latitude, req.Longitude, incidents)
			if err := s.webhookService.Enqueue(context.Background(), task); err != nil {
				s.logger.Warn("failed to enqueue webhook", "error", err)
			}
		}
	}()

	s.logger.Info("location checked",
		"user_id", req.UserID,
		"lat", req.Latitude,
		"lon", req.Longitude,
		"incidents_found", len(incidents))

	return response, nil
}
