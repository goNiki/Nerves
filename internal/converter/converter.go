package converter

import (
	"github.com/goNiki/Nerves/internal/dto"
	"github.com/goNiki/Nerves/internal/models"
	"github.com/google/uuid"
)

// Если уровень опасности не обозначен, то по умолчанию unknown
// если радиус не обозначен, то по умолчанию 100 метров
func CreateIncidentRequestsToModel(req *dto.CreateIncidentRequests) *models.Incident {
	severity := models.SeverityUnknown

	if req.Severity != "" {
		severity = models.Severity(req.Severity)
	}

	radiusMeters := 100
	if req.RadiusMeters > 0 {
		radiusMeters = req.RadiusMeters
	}

	return &models.Incident{
		Title:        req.Title,
		Description:  req.Description,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		RadiusMeters: radiusMeters,
		Severity:     severity,
		Active:       true,
	}
}

func UpdateIncidentRequestsToModel(req *dto.UpdateIncidentRequests) *models.UpdateIncident {
	resp := &models.UpdateIncident{
		ID: uuid.MustParse(req.IncidentID),
	}

	if req.Title != nil {
		resp.Title = req.Title
	}

	if req.Description != nil {
		resp.Description = req.Description
	}

	if req.Latitude != nil {
		resp.Latitude = req.Latitude
	}

	if req.Longitude != nil {
		resp.Longitude = req.Longitude
	}

	if req.RadiusMeters != nil {
		resp.RadiusMeters = req.RadiusMeters
	}

	if req.Severity != nil {
		severity := models.Severity(*req.Severity)
		resp.Severity = &severity
	}

	if req.Active != nil {
		resp.Active = req.Active
	}

	return resp
}

// ApplyUpdate применяет частичное обновление к существующему инциденту
// Возвращает новый Incident с применёнными изменениями
func ApplyUpdateToIncident(existing *models.Incident, update *models.UpdateIncident) *models.Incident {
	result := *existing // копируем
	result.ID = update.ID

	if update.Title != nil {
		result.Title = *update.Title
	}
	if update.Description != nil {
		result.Description = *update.Description
	}
	if update.Latitude != nil {
		result.Latitude = *update.Latitude
	}
	if update.Longitude != nil {
		result.Longitude = *update.Longitude
	}
	if update.RadiusMeters != nil {
		result.RadiusMeters = *update.RadiusMeters
	}
	if update.Severity != nil {
		result.Severity = *update.Severity
	}
	if update.Active != nil {
		result.Active = *update.Active
	}

	return &result
}

func IncidentToResponce(incident *models.Incident) dto.IncidentResponse {
	return dto.IncidentResponse{
		ID:            incident.ID.String(),
		Title:         incident.Title,
		Description:   incident.Description,
		Latitude:      incident.Latitude,
		Longitude:     incident.Longitude,
		RadiusMeters:  incident.RadiusMeters,
		Severity:      string(incident.Severity),
		Active:        incident.Active,
		CreatedAt:     incident.CreateAt,
		UpdatedAt:     incident.UpdatedAt,
		DeactivatedAt: incident.DeactivatedAt,
	}
}

func IncidentListToResponce(list *models.IncidentList) dto.IncidentListResponse {
	items := make([]dto.IncidentResponse, len(list.Items))

	for i, incident := range list.Items {
		items[i] = IncidentToResponce(incident)
	}

	return dto.IncidentListResponse{
		Items: items,
		Total: list.Total,
		Page:  list.Page,
		Limit: list.Limit,
	}

}

func CheckLocationRequestToModel(req *dto.CheckLocationRequest) *models.LocationCheckRequest {
	return &models.LocationCheckRequest{
		UserID:    req.UserID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}
}

func LocationCheckToResponce(resp *models.LocationCheckResponce) *dto.CheckLocationResponce {
	incidents := make([]dto.IncidentResponse, len(resp.Incident))

	for i, incident := range resp.Incident {
		incidents[i] = IncidentToResponce(incident)
	}

	return &dto.CheckLocationResponce{
		Incidents: incidents,
		CheckedAt: resp.CheckedAt,
	}

}
