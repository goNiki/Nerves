package repoModels

import (
	"time"

	"github.com/google/uuid"
)

type IncidentRow struct {
	ID            uuid.UUID
	Title         string
	Description   string
	Latitude      float64
	Longitude     float64
	RadiusMeters  int
	Severity      string
	Active        bool
	CreateAt      time.Time
	UpdatedAt     *time.Time
	DeactivatedAt *time.Time
}

type LocationCheckRow struct {
	ID             uuid.UUID
	UserID         string
	Latitude       float64
	Longitude      float64
	IncidentsFound int
	IncidentIDs    []uuid.UUID
	CheckedAt      time.Time
}

type IncidentUserCountRow struct {
	IncidentID uuid.UUID
	Title      string
	UserCount  int
}
