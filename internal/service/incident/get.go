package incident

import (
	"context"

	"github.com/goNiki/Nerves/internal/models"
	"github.com/google/uuid"
)

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Incident, error) {

	incident, err := s.dataSource.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get incident", "id", id, "error", err)
		return nil, err
	}

	return incident, nil
}

func (s *Service) List(ctx context.Context, page, limit int) (*models.IncidentList, error) {

	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	incident, total, err := s.dataSource.GetListIncident(ctx, offset, limit)
	if err != nil {
		s.logger.Error("failed to list incedints", "error", err)
		return nil, err
	}

	return &models.IncidentList{
		Items: incident,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil

}
