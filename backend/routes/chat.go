package routes

import (
	"backend/controllers"
	"backend/interfaces"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupChatRoutes(router *gin.Engine, db interfaces.DatabaseService) {
	chatController := controllers.NewChatController(db)

	// Chat routes (all protected by auth)
	chatGroup := router.Group("/api/chat")
	chatGroup.Use(middleware.AuthMiddleware())
	{
		// Chat management
		chatGroup.POST("/create", chatController.CreateChat)
		chatGroup.POST("/:chatID/message", chatController.SendMessage)
		chatGroup.POST("/:chatID/message/stream", chatController.SendMessageStream)
		chatGroup.GET("/list", chatController.GetChats)
		chatGroup.GET("/:chatID", chatController.GetChat)
		chatGroup.DELETE("/:chatID", chatController.DeleteChat)

		// OpenRouter configuration
		chatGroup.GET("/models", chatController.GetModels)
		chatGroup.POST("/apikey", chatController.SetAPIKey)
		chatGroup.GET("/apikey/status", chatController.GetAPIKeyStatus)
		chatGroup.GET("/credits", chatController.GetCredits) // Add new endpoint for credits
	}
}
