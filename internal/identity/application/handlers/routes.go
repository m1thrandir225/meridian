package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/meridian/internal/identity/application/services"
	"github.com/m1thrandir225/meridian/pkg/auth"
)

func SetupIdentityRouter(
	service *services.IdentityService,
	tokenVerifier auth.TokenVerifier,
	integrationGrpcURL string,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	handler := NewHTTPHandler(service)
	authHandler := NewAuthHandler(service, tokenVerifier, integrationGrpcURL)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "identity"})
	})

	apiV1 := router.Group("/api/v1/auth")
	{
		apiV1.POST("/register", handler.handleRegisterRequest)
		apiV1.POST("/login", handler.handleLoginRequest)

		apiV1.GET("/validate-token", authHandler.ValidateToken)

		me := apiV1.Group("/me")
		me.Use(AuthenticationMiddleware(tokenVerifier))
		{
			me.GET("", handler.handleGetCurrentUser)
			me.DELETE("", handler.handleDeleteUserRequest)
			me.PUT("/update-profile", handler.handleUpdateCurrentUserRequest)
			me.PUT("/password", handler.handleUpdateUserPasswordRequest)
			me.POST("/refresh-token", handler.handleRefreshTokenRequest)
		}
	}
	log.Println("Identity HTTP Router configured")
	return router
}
