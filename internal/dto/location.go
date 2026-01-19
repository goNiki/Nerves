package dto

import "time"

type CheckLocationRequest struct {
	UserID    string  `json:"user_id" binding:"required,min=1"`
	Latitude  float64 `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" binding:"required,min=-180,max=180"`
}

//уточнить, тут нужна пагинация или нет
type CheckLocationResponce struct {
	Incidents []IncidentResponse `json:"incidents"`
	CheckedAt time.Time          `json:"checked_at"`
}
