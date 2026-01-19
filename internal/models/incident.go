package models

import (
	"time"

	"github.com/google/uuid"
)

type Severity string

const (
	SeverityUnknown  Severity = "unknown"
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "High"
	SeverityCritical Severity = "critical"
)

type Incident struct {
	ID            uuid.UUID  `json:"id"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	Latitude      float64    `json:"latitude"`
	Longitude     float64    `json:"longitude"`
	RadiusMeters  int        `json:"radius_meters"`
	Severity      Severity   `json:"severity"`
	Active        bool       `json:"active"`
	CreateAt      time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	DeactivatedAt *time.Time `json:"deactivated_at,omitempty"`
}

type IncidentWithDistanceInfo struct {
	Incident                 *Incident
	DistanceToCenter         float64 // расстояние до центра инцидента в метрах
	DistanceToIncidentBorder float64 // расстояние до границы радиуса инцидента (отрицательное если внутри)
	IsUserInside             bool    // находится ли пользователь внутри радиуса инцидента
	IsInUserRange            bool    // попадает ли центр инцидента в радиус поиска пользователя
}

type UpdateIncident struct {
	ID           uuid.UUID
	Title        *string
	Description  *string
	Latitude     *float64
	Longitude    *float64
	RadiusMeters *int
	Severity     *Severity
	Active       *bool
}

type IncidentList struct {
	Items []*Incident
	Total int
	Page  int
	Limit int
}

type IncidentUserCount struct {
	IncidentID uuid.UUID `json:"incident_id"`
	Title      string    `json:"title"`
	UserCount  int       `json:"user_count"`
}

type IncidentStatsList struct {
	Items             []IncidentUserCount `json:"items"`
	TimeWindowMinutes int                 `json:"time_window_minutes"`
}
