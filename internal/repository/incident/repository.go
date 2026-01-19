package incident

import "github.com/goNiki/Nerves/internal/infrastructure/database"

type Repository struct {
	db *database.Postgres
}

func NewIncidentRepository(db *database.Postgres) *Repository {
	return &Repository{
		db: db,
	}
}
