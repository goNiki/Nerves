package repoConverter

import (
	"github.com/goNiki/Nerves/internal/models"
	repoModels "github.com/goNiki/Nerves/internal/repository/models"
)

func IncidentToRow(incident *models.Incident) *repoModels.IncidentRow {
	return &repoModels.IncidentRow{
		ID:            incident.ID,
		Title:         incident.Title,
		Description:   incident.Description,
		Latitude:      incident.Latitude,
		Longitude:     incident.Longitude,
		RadiusMeters:  incident.RadiusMeters,
		Severity:      string(incident.Severity),
		Active:        incident.Active,
		CreateAt:      incident.CreateAt,
		UpdatedAt:     incident.UpdatedAt,
		DeactivatedAt: incident.DeactivatedAt,
	}
}

func RowToIncident(row *repoModels.IncidentRow) *models.Incident {
	return &models.Incident{
		ID:            row.ID,
		Title:         row.Title,
		Description:   row.Description,
		Latitude:      row.Latitude,
		Longitude:     row.Longitude,
		RadiusMeters:  row.RadiusMeters,
		Severity:      models.Severity(row.Severity),
		Active:        row.Active,
		CreateAt:      row.CreateAt,
		UpdatedAt:     row.UpdatedAt,
		DeactivatedAt: row.DeactivatedAt,
	}
}

func RowsToIncidents(rows []*repoModels.IncidentRow) []*models.Incident {
	result := make([]*models.Incident, len(rows))

	for i, row := range rows {
		result[i] = RowToIncident(row)
	}

	return result
}

func LocationCheckRowToModels(check *repoModels.LocationCheckRow) *models.LocationCheck {
	return &models.LocationCheck{
		ID:            check.ID,
		UserID:        check.UserID,
		Latitude:      check.Latitude,
		Longtude:      check.Longitude,
		IncidentFound: check.IncidentsFound,
		IncidentIDs:   check.IncidentIDs,
		CheckedAt:     check.CheckedAt,
	}
}

func IncidentUserCountRowToModels(stat *repoModels.IncidentUserCountRow) *models.IncidentUserCount {
	return &models.IncidentUserCount{
		IncidentID: stat.IncidentID,
		Title:      stat.Title,
		UserCount:  stat.UserCount,
	}
}

func IncidentStatsLisToModels(rows []repoModels.IncidentUserCountRow) *models.IncidentStatsList {
	var resp models.IncidentStatsList

	for _, row := range rows {
		resp.Items = append(resp.Items, *IncidentUserCountRowToModels(&row))
	}
	return &resp
}
