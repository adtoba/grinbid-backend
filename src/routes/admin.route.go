package routes

import (
	"github.com/adtoba/grinbid-backend/src/controllers"
	"github.com/adtoba/grinbid-backend/src/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type AdminRouteController struct {
	adminController  controllers.AdminController
	walletController controllers.WalletController
}

func NewAdminRouteController(adminController controllers.AdminController, walletController controllers.WalletController) *AdminRouteController {
	return &AdminRouteController{adminController, walletController}
}

func (rc *AdminRouteController) RegisterRoutes(rg *gin.RouterGroup, redisClient *redis.Client) {
	router := rg.Group("/admin")
	router.GET("/users", middleware.AuthMiddleware(redisClient), middleware.IsAdmin(), rc.adminController.GetAllUsers)
	router.GET("/users/:id", middleware.AuthMiddleware(redisClient), middleware.IsAdmin(), rc.adminController.GetUser)
	router.POST("/users/:id/block", middleware.AuthMiddleware(redisClient), middleware.IsAdmin(), rc.adminController.BlockUser)
	router.POST("/users/:id/unblock", middleware.AuthMiddleware(redisClient), middleware.IsAdmin(), rc.adminController.UnblockUser)
	router.POST("/categories", middleware.AuthMiddleware(redisClient), middleware.IsAdmin(), rc.adminController.CreateCategory)
	router.GET("/wallets/transactions", middleware.AuthMiddleware(redisClient), middleware.IsAdmin(), rc.walletController.GetAllWalletTransactions)
	router.GET("/wallets/transactions/:id", middleware.AuthMiddleware(redisClient), middleware.IsAdmin(), rc.walletController.GetWalletTransactionById)
	router.GET("/wallets/transactions/user/:id", middleware.AuthMiddleware(redisClient), middleware.IsAdmin(), rc.walletController.GetWalletTransactionsByUserId)
}
