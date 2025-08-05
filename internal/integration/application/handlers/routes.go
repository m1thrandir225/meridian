package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/meridian/internal/integration/application/services"
)

func SetupIntegrationRouter(service *services.IntegrationService) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	handler := NewHttpHandler(service)

	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/integrations", handler.handleRegisterIntegration)
		apiV1.DELETE("/integrations", handler.handleRevokeIntegration)
	}
	log.Println("Integration HTTP Router configured")
	return router
}
