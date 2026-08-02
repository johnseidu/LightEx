package bootstrap

import (
	"fmt"
	"log/slog"

	"github.com/light-group/light-ex-backend/internal/platform/config"
	"github.com/light-group/light-ex-backend/internal/platform/database"
	"github.com/light-group/light-ex-backend/internal/platform/logger"
)

type Application struct {
	Config   *config.Config
	Logger   *slog.Logger
	Database *database.Database
}

func Start() (*Application, error) {

	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load configuration: %w", err)
	}

	log := logger.New()

	log.Info("Starting LightEx...")

	db, err := database.New(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("initialize database: %w", err)
	}

	log.Info("PostgreSQL connected successfully.")

	// Run database migrations.
	if err := database.RunMigrations(cfg.Database); err != nil {
		return nil, fmt.Errorf("run database migrations: %w", err)
	}

	log.Info("Database migrations completed successfully.")

	app := &Application{
		Config:   cfg,
		Logger:   log,
		Database: db,
	}

	log.Info("Bootstrap completed successfully.")

	return app, nil
}

func (a *Application) Close() error {

	if a == nil {
		return nil
	}

	if a.Database != nil {
		return a.Database.Close()
	}

	return nil
}
