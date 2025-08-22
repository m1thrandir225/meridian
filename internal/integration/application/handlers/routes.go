package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/meridian/internal/integration/application/services"
	"github.com/m1thrandir225/meridian/pkg/cache"
	"github.com/m1thrandir225/meridian/pkg/logging"
)

func SetupIntegrationRouter(
	service *services.IntegrationService,
	cache *cache.RedisCache,
	messageClient *services.MessagingClient,
	logger *logging.Logger,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(logging.GinLoggingMiddleware(logger))
	router.Use(logging.GinRecoveryMiddleware(logger))

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
