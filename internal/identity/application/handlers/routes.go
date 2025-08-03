package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/meridian/internal/identity/application/services"
	"github.com/m1thrandir225/meridian/pkg/auth"
	"log"
)

func SetupIdentityRouter(service *services.IdentityService, tokenVerifier auth.TokenVerifier) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	handler := NewHTTPHandler(service)

	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/register", handler.handleRegisterRequest)
		apiV1.POST("/login", handler.handleLoginRequest)

		me := apiV1.Group("/me")
		me.Use(AuthenticationMiddleware(tokenVerifier))
		{
			me.GET("", handler.handleGetCurrentUser)
			me.DELETE("", handler.DeleteCurrentUser)
			me.PUT("/update-profile", handler.UpdateCurrentUser)
		}
	}
	log.Println("Identity HTTP Router configured")
	return router
}
