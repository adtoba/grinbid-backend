package controllers

import (
	"net/http"
	"sort"
	"time"

	"github.com/adtoba/grinbid-backend/src/models"
	"github.com/gin-gonic/gin"
	"github.com/pusher/pusher-http-go"
	"gorm.io/gorm"
)

type ChatController struct {
	DB           *gorm.DB
	PusherClient *pusher.Client
}

func NewChatController(db *gorm.DB, pusherClient *pusher.Client) *ChatController {
	return &ChatController{DB: db, PusherClient: pusherClient}
}

// CreateChat godoc
// @Summary Create a new chat
// @Description Create a new chat for a project
// @Tags chats
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param chat body models.CreateChatRequest true "Chat details"
// @Success 200 {object} models.ChatResponse "Chat created successfully"
// @Router /projects/{id}/chat [post]
func (cc *ChatController) CreateChat(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	var body models.CreateChatRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse("Invalid request", nil))
		return
	}

	chat := models.Chat{
		ListingID:    body.ListingID,
		ListingName:  body.ListingName,
		Participants: models.StringArray(body.Participants),
		UserID:       userID,
		IsGroupChat:  body.IsGroupChat,
		GroupName:    body.GroupName,
	}

	result := cc.DB.Create(&chat)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Failed to create chat", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Chat created successfully", chat))
}

// GetChatByID godoc
// @Summary Get a chat by ID
// @Description Get a chat by its ID
// @Tags chats
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Chat ID"
// @Success 200 {object} models.ChatResponse "Chat fetched successfully"
// @Router /projects/{id}/chat/{chat_id} [get]
func (cc *ChatController) GetChatByID(c *gin.Context) {
	chatID := c.Param("id")

	var chat models.Chat
	cc.DB.Where("id = ?", chatID).First(&chat)

	c.JSON(http.StatusOK, models.SuccessResponse("Chat fetched successfully", chat))
}

// GetUserChats godoc
// @Summary Get all chats for a user
// @Description Get all chats where the user is a participant
// @Tags chats
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {array} models.ChatResponse "Chats fetched successfully"
// @Router /projects/{id}/chat [get]
func (cc *ChatController) GetUserChats(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	var chats []models.Chat
	if err := cc.DB.Where("? = ANY(participants)", userID).Find(&chats).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse("Failed to fetch chats", nil))
		return
	}

	var chatResponses []models.ChatResponse
	for _, chat := range chats {
		var members []models.SimpleUserResponse
		for _, participantID := range chat.Participants {
			var member models.User
			if err := cc.DB.Where("id = ?", participantID).First(&member).Error; err == nil {
				var user models.User
				if err := cc.DB.Where("id = ?", member.ID).First(&user).Error; err == nil {
					members = append(members, models.SimpleUserResponse{
						ID:              member.ID,
						Username:        user.Username,
						FullName:        user.FullName,
						Email:           user.Email,
						Phone:           user.Phone,
						Location:        user.Location,
						IsVerified:      user.IsVerified,
						IsBlocked:       user.IsBlocked,
						IsEmailVerified: user.IsEmailVerified,
						IsPhoneVerified: user.IsPhoneVerified,
						CreatedAt:       user.CreatedAt,
						UpdatedAt:       user.UpdatedAt,
					})
				}
			}
		}
		var lastMessage models.Message
		cc.DB.Where("chat_id = ?", chat.ID).Order("created_at DESC").First(&lastMessage)
		chatResponses = append(chatResponses, chat.ToChatResponse(members, lastMessage.ToMessageResponse(models.User{})))
	}

	// Sort chatResponses by last message timestamp and creation time
	sort.Slice(chatResponses, func(i, j int) bool {
		// Get the latest timestamp for each chat (either last message or creation time)
		var iTime, jTime time.Time

		// For chat i
		if !chatResponses[i].LastMessage.CreatedAt.IsZero() {
			iTime = chatResponses[i].LastMessage.CreatedAt
		} else {
			iTime = chatResponses[i].CreatedAt
		}

		// For chat j
		if !chatResponses[j].LastMessage.CreatedAt.IsZero() {
			jTime = chatResponses[j].LastMessage.CreatedAt
		} else {
			jTime = chatResponses[j].CreatedAt
		}

		// Sort by the latest timestamp (most recent first)
		return iTime.After(jTime)
	})

	c.JSON(http.StatusOK, models.SuccessResponse("Chats fetched successfully", chatResponses))
}
