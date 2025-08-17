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
	messageClient *services.MessagingClient,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	handler := NewHttpHandler(service, cache, messageClient)

	router.GET("/health", handler.handleGetHealth)
	router.GET("/metrics", handler.handleGetMetrics)

	apiV1 := router.Group("/api/v1")
	{
		integrations := apiV1.Group("/integrations")
		{
			integrations.POST("", handler.handleRegisterIntegration)
			integrations.DELETE("/:id", handler.handleRevokeIntegration)
			integrations.PUT("/:id", handler.handleUpdateIntegration)
			integrations.GET("", handler.handleListIntegrations)
		}

		webhookGroup := apiV1.Group("/integrations/webhook")
		{
			webhookGroup.POST("/message", handler.handleWebhookMessage)
		}

		callbackGroup := apiV1.Group("/integrations/callback")
		{
			callbackGroup.POST("/message", handler.handleCallbackMessage)
		}
	}
	log.Println("Integration HTTP Router configured")
	return router
}
