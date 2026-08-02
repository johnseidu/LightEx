package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all application routes.
func (s *Server) RegisterRoutes() {

	// Health Check
	s.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "LightEx",
			"version": "0.1.0",
		})
	})
}
