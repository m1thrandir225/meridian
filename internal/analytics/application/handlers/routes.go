package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/meridian/internal/analytics/application/services"
	"github.com/m1thrandir225/meridian/pkg/logging"
)

func SetupAnalyticsRouter(
	analyticsService *services.AnalyticsService,
	logger *logging.Logger,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(logging.GinLoggingMiddleware(logger))
	router.Use(logging.GinRecoveryMiddleware(logger))

	handler := NewHTTPHandler(analyticsService, logger)

	router.GET("/health", handler.handleGetHealth)

	apiV1 := router.Group("/api/v1/analytics")
	apiV1.Use(AdminMiddleware())
	{
		apiV1.GET("/dashboard", handler.handleGetDashboard)
		apiV1.GET("/user-growth", handler.handleGetUserGrowth)
		apiV1.GET("/message-volume", handler.handleGetMessageVolume)
		apiV1.GET("/channel-activity", handler.handleGetChannelActivity)
		apiV1.GET("/top-users", handler.handleGetTopUsers)
		apiV1.GET("/reaction-usage", handler.handleGetReactionUsage)
	}

	log.Println("Analytics HTTP Router configured")
	return router
}
