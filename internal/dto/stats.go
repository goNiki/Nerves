package dto

type IncidentStatsResponse struct {
	IncidentID string `json:"incident_id"`
	Title      string `json:"title"`
	UserCount  int    `json:"user_count"`
}

type IncidentStatsListResponse struct {
	Items             []IncidentStatsResponse `json:"items"`
	TimeWindowMinutes int                     `json:"time_window_minutes"`
}
