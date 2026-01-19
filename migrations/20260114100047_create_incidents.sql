-- +goose Up
-- +goose StatementBegin

CREATE TABLE incidents(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    radius_meters DOUBLE PRECISION NOT NULL DEFAULT 1000.0,
    severity VARCHAR(20) NOT NULL DEFAULT 'unknown',
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);
-- Индекс для быстрого поиска активных инцидентов
CREATE INDEX idx_incidents_active ON incidents(active) WHERE active = true;

-- Индекс для поиска по координатам
CREATE INDEX idx_incidents_coords ON incidents(latitude, longitude) WHERE active = true;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_incidents_coords;
DROP INDEX IF EXISTS idx_incidents_active;
DROP TABLE IF EXISTS incidents;

-- +goose StatementEnd
