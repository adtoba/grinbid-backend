package routes

import (
	"github.com/adtoba/grinbid-backend/src/controllers"
	"github.com/adtoba/grinbid-backend/src/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type ChatRouteController struct {
	chatController controllers.ChatController
}

func NewChatRouteController(chatController controllers.ChatController) *ChatRouteController {
	return &ChatRouteController{chatController: chatController}
}

func (cr *ChatRouteController) RegisterRoutes(rg *gin.RouterGroup, redisClient *redis.Client) {
	router := rg.Group("/chat", middleware.AuthMiddleware(redisClient))
	router.POST("", cr.chatController.CreateChat)
	router.GET("/:chat_id", cr.chatController.GetChatByID)
	router.GET("", cr.chatController.GetUserChats)
}
