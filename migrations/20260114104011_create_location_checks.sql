-- +goose Up
-- +goose StatementBegin
CREATE TABLE location_checks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    incidents_found INTEGER NOT NULL DEFAULT 0,
    incident_ids UUID[] DEFAULT '{}',
    checked_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- Индекс для поиска по пользователю
CREATE INDEX idx_location_checks_user ON location_checks(user_id);

-- Индекс для поиска по времени (для статистики)
CREATE INDEX idx_location_checks_time ON location_checks(checked_at);

-- GIN индекс для поиска по массиву инцидентов
CREATE INDEX idx_location_checks_incidents ON location_checks USING GIN(incident_ids);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_location_checks_incidents;
DROP INDEX IF EXISTS idx_location_checks_time;
DROP INDEX IF EXISTS idx_location_checks_user;
DROP TABLE IF EXISTS location_checks;
-- +goose StatementEnd
