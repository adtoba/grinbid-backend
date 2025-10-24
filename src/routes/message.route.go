package routes

import (
	"github.com/adtoba/grinbid-backend/src/controllers"
	"github.com/adtoba/grinbid-backend/src/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type MessageRouteController struct {
	messageController controllers.MessageController
}

func NewMessageRouteController(messageController controllers.MessageController) *MessageRouteController {
	return &MessageRouteController{messageController: messageController}
}

func (mr *MessageRouteController) RegisterRoutes(rg *gin.RouterGroup, redisClient *redis.Client) {
	router := rg.Group("/chat/:chat_id/messages", middleware.AuthMiddleware(redisClient))
	router.POST("", mr.messageController.SendMessage)
	router.GET("", mr.messageController.GetMessages)
	router.POST("/:message_id/seen", mr.messageController.MarkMessageAsSeen)
}
