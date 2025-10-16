package routes

import (
	"github.com/adtoba/grinbid-backend/src/controllers"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type WebhookRouteController struct {
	webhookController controllers.WebhooksController
}

func NewWebhookRouteController(webhookController controllers.WebhooksController) *WebhookRouteController {
	return &WebhookRouteController{webhookController}
}

func (rc *WebhookRouteController) RegisterRoutes(rg *gin.RouterGroup, redisClient *redis.Client) {
	rg.POST("/paystack", rc.webhookController.PaystackWebhook)
}
