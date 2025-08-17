package handlers

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	httpHandler *HTTPHandler,
	wsHandler *WebSocketHandler,
) {
	router.GET("/health", httpHandler.handleGetHealth)
	router.GET("/metrics", httpHandler.handleGetMetrics)

	apiV1 := router.Group("/api/v1/messages")
	{
		apiV1.GET("/ws", func(c *gin.Context) {
			wsHandler.HandleWebSocket(c)
		})

		channelsGroup := apiV1.Group("/channels")
		{
			channelsGroup.GET("/", httpHandler.handleGetUserChannels)
			channelsGroup.POST("/", httpHandler.handleCreateChannel)
			channelsGroup.GET("/:channelId", httpHandler.handleGetChannel)
			channelsGroup.POST("/:channelId/join", httpHandler.handleJoinChannel)
			channelsGroup.PUT("/:channelId/archive", httpHandler.handleArchiveChannel)
			channelsGroup.PUT("/:channelId/unarchive", httpHandler.handleUnarchiveChannel)
			channelsGroup.POST("/:channelId/bots", httpHandler.handleAddBotToChannel)

			messagesGroup := channelsGroup.Group("/:channelId/messages")
			{
				messagesGroup.GET("", httpHandler.handleGetMessages)
				messagesGroup.POST("", httpHandler.handleSendMessage)

				reactionsGroup := messagesGroup.Group("/:messageId/reactions")
				{
					reactionsGroup.PUT("/", httpHandler.handleAddReaction)
					reactionsGroup.DELETE("/", httpHandler.handleRemoveReaction)
				}
			}
		}
	}
}
