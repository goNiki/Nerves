package di

import (
	"github.com/gin-gonic/gin"
	"github.com/goNiki/Nerves/internal/handlers/health"
	"github.com/goNiki/Nerves/internal/handlers/incident"
	"github.com/goNiki/Nerves/internal/handlers/location"
	"github.com/goNiki/Nerves/internal/handlers/stats"
	"github.com/goNiki/Nerves/internal/infrastructure/httpserver/middleware"
)

func (c *Container) SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	logger := c.Logger.Log

	r.Use(middleware.Recovery(logger))
	r.Use(middleware.Logger(logger))

	api := r.Group("/api/v1")

	// публичные роуты
	location.RegisterRoutes(api, c.LocationHandler)
	health.RegisterRoutes(api, c.HealthHandler)

	// защищённые роуты
	protected := api.Group("")
	protected.Use(middleware.APIKeyAuth(c.Config.Api.APIKey()))

	incident.RegisterRoutes(protected, c.IncidentHandler)
	stats.RegisterRoutes(protected, c.StatsHandler)

	logger.Info("routes registered")

	return r
}
