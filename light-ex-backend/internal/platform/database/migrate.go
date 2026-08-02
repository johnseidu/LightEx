package database

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/light-group/light-ex-backend/internal/platform/config"
)

func RunMigrations(cfg config.DatabaseConfig) error {

	m, err := migrate.New(
		"file://migrations",
		cfg.MigrationURL(),
	)
	if err != nil {
		return fmt.Errorf("create migration instance: %w", err)
	}

	defer func() {
		_, _ = m.Close()
	}()

	err = m.Up()

	if err != nil {

		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}

		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
