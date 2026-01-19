package models

type HealthStatus struct {
	Status     string
	Components map[string]string
}
