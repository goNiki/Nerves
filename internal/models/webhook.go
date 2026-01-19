package models

import (
	"time"

	"github.com/google/uuid"
)

type WebhookTask struct {
	ID        uuid.UUID   `json:"id"`
	UserID    string      `json:"user_id"`
	Latitude  float64     `json:"latitude"`
	Longitude float64     `json:"longitude"`
	Incidents []*Incident `json:"incidents"`
	Attempt   int         `json:"attempt"`
	CreatedAt time.Time   `json:"created_at"`
}

func NewWebhookTask(userID string, lat, lon float64, incidents []*Incident) *WebhookTask {
	return &WebhookTask{
		ID:        uuid.New(),
		UserID:    userID,
		Latitude:  lat,
		Longitude: lon,
		Incidents: incidents,
		Attempt:   0,
		CreatedAt: time.Now(),
	}
}
