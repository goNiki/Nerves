package errorapp

import "errors"

//Внутренние ошибки связанные с базой данных

var (
	ErrCreateIncident     = errors.New("failed to create incident")
	ErrCountIncidents     = errors.New("failed to count incidents")
	ErrListIncidents      = errors.New("failed to list incidents")
	ErrScanIncident       = errors.New("failed to scan incident")
	ErrGetIncident        = errors.New("failed to get incident")
	ErrUpdateIncident     = errors.New("failed to update incident")
	ErrDeactivateIncident = errors.New("failed to deactivate incident")
	ErrIncidentNotFound   = errors.New("incident not found")
)

var (
	ErrCreateLocationCheck = errors.New("failed to create location check")
	ErrGetIncidentStats    = errors.New("failed to get incidents stats")
	ErrScanIcidentStats    = errors.New("failed to scan stats row")
	ErrGetFromCache        = errors.New("failed to get from cache")
	ErrSetToCache          = errors.New("failed to set to cache")
)

var (
	ErrGetCacheAllActiveIncidents = errors.New("failed to get cache all actives incidents")
	ErrGetIncidentsByIDs          = errors.New("failde to get cashe incidents by ids ")
	ErrGetCacheActiveIncident     = errors.New("failed to get cache actives incident")
	ErrUnmarshallCashIncident     = errors.New("failed to unmarshal cash incident")
	ErrMarshallCashIncident       = errors.New("failed to marshall cash incident")
	ErrDeleteCashIncident         = errors.New("failed to delete cash incident")
	ErrMaxRadiusNotFound          = errors.New("max radius incidents is not found ")
	ErrGeoSearchIncidentIDs       = errors.New("failed to ger search incidents ids")
)

var (
	ErrEnqueueWebhook    = errors.New("failed to marshall enqueuewebhook")
	ErrDequeueWebhook    = errors.New("failed to marshall dequeuewebhook")
	ErrDeleteWebhookTask = errors.New("failed to delete webhook task")
)
