package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/meridian/internal/identity/application/services"
	"github.com/m1thrandir225/meridian/pkg/auth"
	"github.com/m1thrandir225/meridian/pkg/cache"
	"github.com/m1thrandir225/meridian/pkg/logging"
)

func SetupIdentityRouter(
	service *services.IdentityService,
	cache *cache.RedisCache,
	tokenVerifier auth.TokenVerifier,
	integrationGrpcURL string,
	logger *logging.Logger,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(logging.GinLoggingMiddleware(logger))
	router.Use(logging.GinRecoveryMiddleware(logger))

	handler := NewHTTPHandler(service, cache)
	authHandler := NewAuthHandler(service, cache, tokenVerifier, integrationGrpcURL)

	router.GET("/health", handler.handleGetHealth)
	router.GET("/metrics", handler.handleGetMetrics)

	apiV1 := router.Group("/api/v1/auth")
	{
		apiV1.POST("/register", handler.handleRegisterRequest)
		apiV1.POST("/login", handler.handleLoginRequest)

		apiV1.GET("/validate-token", authHandler.ValidateToken)
		apiV1.POST("/refresh-token", handler.handleRefreshTokenRequest)

		me := apiV1.Group("/me")
		me.Use(AuthenticationMiddleware(tokenVerifier))
		{
			me.GET("", handler.handleGetCurrentUser)
			me.DELETE("", handler.handleDeleteUserRequest)
			me.PUT("/update-profile", handler.handleUpdateCurrentUserRequest)
			me.PUT("/password", handler.handleUpdateUserPasswordRequest)
		}
	}
	log.Println("Identity HTTP Router configured")
	return router
}
