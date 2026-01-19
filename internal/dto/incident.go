package dto

import (
	"time"
)

type CreateIncidentRequests struct {
	Title        string  `json:"title" binging:"required, min=1,max=255"`
	Description  string  `json:"description"`
	Latitude     float64 `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude    float64 `json:"longitude" binding:"required,min=-180,max=180"`
	RadiusMeters int     `json:"radius_meters" binding:"omitempty,min=1,max=100000"`
	Severity     string  `json:"severity" binding:"omitempty,oneof=low medium high critical"`
}

type UpdateIncidentRequests struct {
	IncidentID   string   `json:"incident_id" binding:"required,uuid"`
	Title        *string  `json:"title" binding:"omitempty,min=1,max=255"`
	Description  *string  `json:"description" binding:"omitempty"`
	Latitude     *float64 `json:"latitude" binding:"omitempty,min=-90,max=90"`
	Longitude    *float64 `json:"longitude" binding:"omitempty,min=-180,max=180"`
	RadiusMeters *int     `json:"radius_meters" binding:"omitempty,min=1,max=100000"`
	Severity     *string  `json:"severity" binding:"omitempty,oneof=low medium high critical"`
	Active       *bool    `json:"active" binding:"omitempty"`
}

type ListIncidentQuery struct {
	Page  int `json:"page" binding:"omitempty,min=1"`
	Limit int `json:"limit" binding:"omitempty,min=1,max=100"`
}

type IncidentResponse struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	Latitude      float64    `json:"latitude"`
	Longitude     float64    `json:"longitude"`
	RadiusMeters  int        `json:"radius_meters"`
	Severity      string     `json:"severity"`
	Active        bool       `json:"active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	DeactivatedAt *time.Time `json:"deactivated_at,omitempty"`
}
type IncidentListResponse struct {
	Items []IncidentResponse `json:"items"`
	Total int                `json:"total"`
	Page  int                `json:"page"`
	Limit int                `json:"limit"`
}

type IncidentWithDistance struct {
	Incident                 *IncidentResponse
	DistanceToCenter         float64 `json:"distance_to_center"`
	DistanceToIncidentBorder float64 `json:"distance_to_incident_border"`
	IsUserInside             float64 `json:"is_user_inside"`
	IsInUserRange            float64 `json:"is_in_user_range"`
}

type CheckLocationWithDistanceResponse struct {
	Incidents []IncidentWithDistance `json:"incidents"`
	CheckedAt time.Time              `json:"checked_at"`
}
