package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/light-group/light-ex-backend/internal/platform/bootstrap"
	"github.com/light-group/light-ex-backend/internal/platform/server"
)

func main() {

	app, err := bootstrap.Start()
	if err != nil {
		log.Fatalf("bootstrap failed: %v", err)
	}

	httpServer := server.New(app.Config, app.Logger)

	RegisterRoutes(
		httpServer.Engine(),
		app,
	)

	go func() {
		if err := httpServer.Start(); err != nil {
			log.Fatalf("server stopped: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	app.Logger.Info("Shutdown signal received.")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		app.Logger.Error("Failed to shutdown HTTP server", "error", err)
	}

	if err := app.Close(); err != nil {
		app.Logger.Error("Failed to close application resources", "error", err)
	}

	app.Logger.Info("LightEx shutdown completed.")
}
