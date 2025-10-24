package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID          string      `json:"id" gorm:"primaryKey"`
	ChatID      string      `json:"chat_id"`
	SenderID    string      `json:"sender_id"`
	ListingID   string      `json:"listing_id"`
	Content     string      `json:"content"`
	MessageType string      `json:"message_type"`
	Attachments StringArray `json:"attachments" gorm:"type:text[]"`
	IsRead      bool        `json:"is_read" gorm:"default:false"`
	IsDeleted   bool        `json:"is_deleted" gorm:"default:false"`
	IsEdited    bool        `json:"is_edited" gorm:"default:false"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type MessageSeen struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	MessageID string    `json:"message_id"`
	UserID    string    `json:"user_id"`
	SeenAt    time.Time `json:"seen_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (message *Message) BeforeCreate(tx *gorm.DB) (err error) {
	message.ID = uuid.New().String()
	return
}

func (ms *MessageSeen) BeforeCreate(tx *gorm.DB) (err error) {
	ms.ID = uuid.New().String()
	ms.SeenAt = time.Now()
	return
}

func (message *Message) ToMessageResponse(sender User) MessageResponse {
	return MessageResponse{
		ID:          message.ID,
		ChatID:      message.ChatID,
		SenderID:    message.SenderID,
		Sender:      sender.ToUserResponse(),
		MessageType: message.MessageType,
		Content:     message.Content,
		Attachments: message.Attachments,
		IsRead:      message.IsRead,
		IsDeleted:   message.IsDeleted,
		IsEdited:    message.IsEdited,
		CreatedAt:   message.CreatedAt,
		UpdatedAt:   message.UpdatedAt,
	}
}

type MessageResponse struct {
	ID          string       `json:"id"`
	ChatID      string       `json:"chat_id"`
	SenderID    string       `json:"sender_id"`
	Sender      UserResponse `json:"sender"`
	MessageType string       `json:"message_type"`
	Attachments StringArray  `json:"attachments"`
	Content     string       `json:"content"`
	IsRead      bool         `json:"is_read"`
	IsDeleted   bool         `json:"is_deleted"`
	IsEdited    bool         `json:"is_edited"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type CreateMessageRequest struct {
	ListingID   string   `json:"listing_id"`
	Content     string   `json:"content" binding:"required"`
	Attachments []string `json:"attachments" gorm:"type:text[]"`
}

type GetMessagesResponse struct {
	ChatID    string            `json:"chat_id"`
	ListingID string            `json:"listing_id"`
	Messages  []MessageResponse `json:"messages"`
}
