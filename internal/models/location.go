package models

import (
	"time"

	"github.com/google/uuid"
)

const DefaultMaxRadius = 10000

type LocationCheck struct {
	ID            uuid.UUID   `json:"id"`
	UserID        string      `json:"user_id"`
	Latitude      float64     `json:"latitude"`
	Longtude      float64     `json:"longitude"`
	IncidentFound int         `json:"incidents_found"`
	IncidentIDs   []uuid.UUID `json:"incident_ids"`
	CheckedAt     time.Time   `json:"checked_at"`
}

type LocationCheckRequest struct {
	UserID    string  `json:"user_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type LocationCheckResponce struct {
	Incident  []*Incident `json:"incidents"`
	CheckedAt time.Time   `json:"checked_at"`
}
