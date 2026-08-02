package auth

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/light-group/light-ex-backend/internal/modules/auth/handler"
	"github.com/light-group/light-ex-backend/internal/modules/auth/repository/postgres"
	"github.com/light-group/light-ex-backend/internal/modules/auth/service"

	"github.com/light-group/light-ex-backend/internal/platform/config"
	jwtservice "github.com/light-group/light-ex-backend/internal/platform/security/jwt"
)

type Module struct {
	handler *handler.Handler
}

// New creates the complete authentication module.
func New(
	cfg *config.Config,
	db *sql.DB,
) *Module {

	// Repository
	repository := postgres.New(db)

	// Password service
	passwordService := service.NewPasswordService()

	// JWT service
	jwt := jwtservice.New(cfg.JWT)

	// Authentication service
	authService := service.New(
		repository,
		passwordService,
		jwt,
	)

	// HTTP handler
	httpHandler := handler.New(authService)

	return &Module{
		handler: httpHandler,
	}
}

// RegisterRoutes registers all authentication routes.
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	m.handler.RegisterRoutes(router)
}