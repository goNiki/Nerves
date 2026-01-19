package location

import "github.com/goNiki/Nerves/internal/infrastructure/database"

type Repository struct {
	db *database.Postgres
}

func NewLocationRepository(db *database.Postgres) *Repository {
	return &Repository{
		db: db,
	}
}
