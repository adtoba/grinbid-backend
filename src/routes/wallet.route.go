package routes

import (
	"github.com/adtoba/grinbid-backend/src/controllers"
	"github.com/adtoba/grinbid-backend/src/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type WalletRouteController struct {
	walletController controllers.WalletController
}

func NewWalletRouteController(walletController controllers.WalletController) *WalletRouteController {
	return &WalletRouteController{walletController}
}

func (rc *WalletRouteController) RegisterRoutes(rg *gin.RouterGroup, redisClient *redis.Client) {
	router := rg.Group("/wallet")
	router.GET("/", middleware.AuthMiddleware(redisClient), rc.walletController.GetWallet)
	router.POST("/purchase", middleware.AuthMiddleware(redisClient), rc.walletController.PurchaseFromWallet)
	router.GET("/transactions", middleware.AuthMiddleware(redisClient), rc.walletController.GetWalletTransactions)
	router.POST("/initialize-transaction", middleware.AuthMiddleware(redisClient), rc.walletController.InitializeTransaction)
}
