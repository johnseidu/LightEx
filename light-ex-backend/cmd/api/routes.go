package main

import (
	"github.com/gin-gonic/gin"
	auth "github.com/light-group/light-ex-backend/internal/modules/auth"
	"github.com/light-group/light-ex-backend/internal/platform/bootstrap"
)

// RegisterRoutes registers all API routes.
func RegisterRoutes(
	router *gin.Engine,
	app *bootstrap.Application) {

	// API root
	api := router.Group("/api")

	// Version 1
	v1 := api.Group("/v1")

	// Health
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "LightEx API",
			"version": "v1",
		})
	})

	// Authentication
	authModule := auth.New(
		app.Config,
		app.Database.SQL(),
	)

	authModule.RegisterRoutes(v1)
}
