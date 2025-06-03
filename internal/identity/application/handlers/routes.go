package handlers

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine, handler *HTTPHandler) {
	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/register", handler.Register)
		apiV1.POST("/login", handler.Login)

		privateRoutes := apiV1.Group("/").Use(AuthenticationMiddleware())
		{
			privateRoutes.GET("/me", handler.GetCurrentUser)
			privateRoutes.DELETE("/me", handler.DeleteCurrentUser)
			privateRoutes.PUT("/update-profile", handler.UpdateCurrentUser)
		}

	}
}
