package controllers

import (
	"log"
	"net/http"

	"github.com/adtoba/grinbid-backend/src/models"
	"github.com/gin-gonic/gin"
	"github.com/pusher/pusher-http-go"
	"gorm.io/gorm"
)

type MessageController struct {
	DB           *gorm.DB
	PusherClient *pusher.Client
}

func NewMessageController(db *gorm.DB, pusherClient *pusher.Client) *MessageController {
	return &MessageController{DB: db, PusherClient: pusherClient}
}

// SendMessage godoc
// @Summary Send a message to a chat
// @Description Send a message to a chat
// @Tags messages
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param chat_id path string true "Chat ID"
// @Param message body models.CreateMessageRequest true "Message"
// @Success 201 {object} models.MessageResponse "Message sent successfully"
// @Router /projects/{id}/chat/{chat_id}/messages [post]
func (mc *MessageController) SendMessage(c *gin.Context) {
	chatID := c.Param("chat_id")
	userID := c.MustGet("user_id")

	var payload *models.CreateMessageRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Invalid request", nil))
		return
	}

	var message models.Message
	message.ChatID = chatID
	message.SenderID = userID.(string)
	message.Content = payload.Content
	message.ListingID = payload.ListingID
	message.MessageType = "text"
	message.Attachments = payload.Attachments

	mc.DB.Create(&message)

	go func() {
		mc.PusherClient.Trigger("grinbid-chat-"+payload.ListingID+"-"+chatID, "new-message", message)
		log.Println("Pusher message sent")
	}()

	c.JSON(http.StatusCreated, models.SuccessResponse("Message sent successfully", message))
}

// GetMessages godoc
// @Summary Get messages for a chat
// @Description Get messages for a chat
// @Tags messages
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param chat_id path string true "Chat ID"
// @Success 200 {object} models.GetMessagesResponse "Messages fetched successfully"
// @Router /projects/{id}/chat/{chat_id}/messages [get]
func (mc *MessageController) GetMessages(c *gin.Context) {
	chatID := c.Param("chat_id")
	userID := c.MustGet("user_id").(string)

	var messages []models.Message
	if err := mc.DB.Where("chat_id = ?", chatID).Find(&messages).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, models.ErrorResponse("Messages not found", nil))
		return
	}

	var chat models.Chat
	if err := mc.DB.Where("id = ?", chatID).First(&chat).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, models.ErrorResponse("Chat not found", nil))
		return
	}

	// Mark messages as seen
	for _, message := range messages {
		// Skip if message is from the current user
		if message.SenderID == userID {
			continue
		}

		// Check if message is already seen
		var existingSeen models.MessageSeen
		if err := mc.DB.Where("message_id = ? AND user_id = ?", message.ID, userID).First(&existingSeen).Error; err != nil {
			// Message not seen yet, create seen record
			messageSeen := models.MessageSeen{
				MessageID: message.ID,
				UserID:    userID,
			}
			if err := mc.DB.Create(&messageSeen).Error; err != nil {
				log.Printf("Failed to mark message as seen: %v", err)
			}
		}
	}

	var messagesResponse []models.MessageResponse
	for _, message := range messages {
		var sender models.User
		mc.DB.Where("id = ?", message.SenderID).First(&sender)

		// Get seen information for this message
		var seenBy []models.MessageSeen
		mc.DB.Where("message_id = ?", message.ID).Find(&seenBy)

		// Check if current user has seen this message
		isRead := false
		for _, seen := range seenBy {
			if seen.UserID == userID {
				isRead = true
				break
			}
		}

		messageResponse := message.ToMessageResponse(sender)
		messageResponse.IsRead = isRead
		messagesResponse = append(messagesResponse, messageResponse)
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Messages fetched successfully", models.GetMessagesResponse{
		ChatID:    chatID,
		ListingID: chat.ListingID,
		Messages:  messagesResponse,
	}))
}

// MarkMessageAsSeen godoc
// @Summary Mark a message as seen
// @Description Mark a specific message as seen by the current user
// @Tags messages
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param chat_id path string true "Chat ID"
// @Param message_id path string true "Message ID"
// @Success 200 {object} models.MessageResponse "Message marked as seen"
// @Router /projects/{id}/chat/{chat_id}/messages/{message_id}/seen [post]
func (mc *MessageController) MarkMessageAsSeen(c *gin.Context) {
	chatID := c.Param("chat_id")
	messageID := c.Param("message_id")
	userID := c.MustGet("user_id").(string)

	// Check if message exists and belongs to the chat
	var message models.Message
	if err := mc.DB.Where("id = ? AND chat_id = ?", messageID, chatID).First(&message).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, models.ErrorResponse("Message not found", nil))
		return
	}

	// Skip if message is from the current user
	if message.SenderID == userID {
		c.JSON(http.StatusOK, models.SuccessResponse("Message is from current user", nil))
		return
	}

	// Check if message is already seen
	var existingSeen models.MessageSeen
	if err := mc.DB.Where("message_id = ? AND user_id = ?", messageID, userID).First(&existingSeen).Error; err != nil {
		// Message not seen yet, create seen record
		messageSeen := models.MessageSeen{
			MessageID: messageID,
			UserID:    userID,
		}
		if err := mc.DB.Create(&messageSeen).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Failed to mark message as seen", nil))
			return
		}
	}

	// Get sender information
	var sender models.User
	mc.DB.Where("id = ?", message.SenderID).First(&sender)

	// Get seen information
	var seenBy []models.MessageSeen
	mc.DB.Where("message_id = ?", messageID).Find(&seenBy)

	// Check if current user has seen this message
	isRead := false
	for _, seen := range seenBy {
		if seen.UserID == userID {
			isRead = true
			break
		}
	}

	messageResponse := message.ToMessageResponse(sender)
	messageResponse.IsRead = isRead

	c.JSON(http.StatusOK, models.SuccessResponse("Message marked as seen", messageResponse))
}
