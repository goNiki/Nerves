package di

import (
	"fmt"

	"github.com/goNiki/Nerves/internal/infrastructure/migrator"
)

func (c *Container) RunMigrations() error {
	m := migrator.NewMigrator(c.DB.DB, "./migrations")

	if err := m.Up(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	c.Logger.Log.Info("migrations completed")
	return nil
}
