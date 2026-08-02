package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/light-group/light-ex-backend/internal/platform/config"
	"github.com/light-group/light-ex-backend/internal/platform/middleware"
)

type Server struct {
	engine *gin.Engine
	http   *http.Server
	logger *slog.Logger
}

func New(cfg *config.Config, log *slog.Logger) *Server {

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()

	engine.Use(
		middleware.Logger(log),
		middleware.Recovery(log),
		middleware.CORS(),
	)

	srv := &Server{
		engine: engine,
		logger: log,
		http: &http.Server{
			Addr:              fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
			Handler:           engine,
			ReadHeaderTimeout: 10 * time.Second,
		},
	}

	return srv
}

func (s *Server) Engine() *gin.Engine {
	return s.engine
}

func (s *Server) Start() error {

	s.logger.Info(
		"Starting HTTP server...",
		slog.String("address", s.http.Addr),
	)

	return s.http.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {

	s.logger.Info("Shutting down HTTP server...")

	if err := s.http.Shutdown(ctx); err != nil {
		return err
	}

	s.logger.Info("HTTP server stopped successfully.")

	return nil
}
