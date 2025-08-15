package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, httpHandler *HTTPHandler, wsHandler *WebSocketHandler) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "messaging"})
	})

	apiV1 := router.Group("/api/v1/messages")
	{
		apiV1.GET("/ws", func(c *gin.Context) {
			wsHandler.HandleWebSocket(c)
		})

		channelsGroup := apiV1.Group("/channels")
		{
			channelsGroup.GET("/", httpHandler.GetUserChannels)
			channelsGroup.POST("/", httpHandler.CreateChannel)
			channelsGroup.GET("/:channelId", httpHandler.GetChannel)
			channelsGroup.POST("/:channelId/join", httpHandler.JoinChannel)
			channelsGroup.PUT("/:channelId/archive", httpHandler.ArchiveChannel)
			channelsGroup.PUT("/:channelId/unarchive", httpHandler.UnarchiveChannel)
			channelsGroup.POST("/:channelId/bots", httpHandler.AddBotToChannel)

			messagesGroup := channelsGroup.Group("/:channelId/messages")
			{
				messagesGroup.GET("/", httpHandler.GetMessages)
				messagesGroup.POST("/", httpHandler.SendMessage)

				reactionsGroup := messagesGroup.Group("/:messageId/reactions")
				{
					reactionsGroup.PUT("/", httpHandler.AddReaction)
					reactionsGroup.DELETE("/", httpHandler.RemoveReaction)
				}
			}
		}
	}
}
