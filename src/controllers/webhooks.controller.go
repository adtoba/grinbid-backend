package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/adtoba/grinbid-backend/src/models"
	"github.com/adtoba/grinbid-backend/src/services"
	"github.com/adtoba/grinbid-backend/src/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WebhooksController struct {
	DB            *gorm.DB
	WalletService *services.WalletService
}

func NewWebhooksController(db *gorm.DB, walletService *services.WalletService) *WebhooksController {
	return &WebhooksController{DB: db, WalletService: walletService}
}

func (wc *WebhooksController) PaystackWebhook(c *gin.Context) {
	var request map[string]interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("invalid request", err.Error()))
		return
	}

	fmt.Println(request)

	var transaction models.Transaction

	if request["event"] == "charge.success" {
		data := request["data"].(map[string]interface{})
		metadata := data["metadata"].(map[string]interface{})
		transactionRef := data["reference"].(string)

		transaction.Status = data["status"].(string)
		transaction.Amount = utils.FromKobo(data["amount"].(float64))
		transaction.Type = metadata["type"].(string)
		transaction.SenderID = metadata["sender_id"].(string)
		transaction.TransactionRef = transactionRef

		if metadata["type"] == "topup" {
			transaction.SenderWalletID = metadata["sender_wallet_id"].(string)
			transaction.PaymentMethod = "paystack"
			transaction.ItemName = "Wallet Topup"
			wc.WalletService.TopupWallet(c, transaction)
		} else {
			transaction.ItemName = metadata["item_name"].(string)
			transaction.ListingID = metadata["listing_id"].(string)
			transaction.ReceiverID = metadata["receiver_id"].(string)
			transaction.SenderWalletID = metadata["sender_wallet_id"].(string)
			transaction.ReceiverWalletID = metadata["receiver_wallet_id"].(string)
			dtt, err := strconv.Atoi(metadata["dtt"].(string))
			if err != nil {
				c.JSON(http.StatusBadRequest, models.ErrorResponse("invalid dtt", err.Error()))
				return
			}
			transaction.Dtt = dtt
			wc.WalletService.PurchaseFromPaystack(c, transaction)
		}

	} else if request["event"] == "transfer.success" {
		data := request["data"].(map[string]interface{})
		recipient := data["recipient"].(map[string]interface{})
		metadata := recipient["metadata"].(map[string]interface{})
		transactionRef := data["reference"].(string)

		transaction.Status = data["status"].(string)
		transaction.Amount = utils.FromKobo(data["amount"].(float64))
		transaction.Type = metadata["type"].(string)
		transaction.TransactionRef = transactionRef
		transaction.SenderID = metadata["sender_id"].(string)
		wc.WalletService.WithdrawFromWallet(c, transaction)
	}

	wc.DB.Create(&transaction)

	// Send a status 200 response to paystack
	c.Status(http.StatusOK)
	c.JSON(http.StatusOK, models.SuccessResponse("Webhook received", request))
}
