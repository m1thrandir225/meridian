package handlers

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine, httpHandler *HTTPHandler) {
	apiV1 := router.Group("/api/v1")
	{
		channelsGroup := apiV1.Group("/channels")
		{
			channelsGroup.POST("/", httpHandler.CreateChannel)
			channelsGroup.GET("/:channelId", httpHandler.GetChannel)
			channelsGroup.POST("/:channelId/join", httpHandler.JoinChannel)
			channelsGroup.PUT("/:channelId/archive", httpHandler.ArchiveChannel)
			channelsGroup.PUT("/:channelId/unarchive", httpHandler.UnarchiveChannel)
			messagesGroup := channelsGroup.Group("/:channelId/messages")
			{
				messagesGroup.POST("/", httpHandler.SendMessage)

				reactionsGroup := messagesGroup.Group("/:messageId/reactions")
				{
					reactionsGroup.PUT("", httpHandler.AddReaction)
					reactionsGroup.DELETE("", httpHandler.RemoveReaction)
				}
			}
		}
	}
}
