package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/meridian/internal/integration/application/services"
	"github.com/m1thrandir225/meridian/pkg/cache"
)

func SetupIntegrationRouter(
	service *services.IntegrationService,
	cache *cache.RedisCache,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	handler := NewHttpHandler(service, cache)

	router.GET("/health", handler.handleGetHealth)
	router.GET("/metrics", handler.handleGetMetrics)

	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/integrations", handler.handleRegisterIntegration)
		apiV1.DELETE("/integrations", handler.handleRevokeIntegration)
	}
	log.Println("Integration HTTP Router configured")
	return router
}
