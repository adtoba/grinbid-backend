package models

import (
	"database/sql/driver"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}

	// Convert to PostgreSQL array literal format
	result := "{"
	for i, v := range a {
		if i > 0 {
			result += ","
		}
		result += "\"" + v + "\""
	}
	result += "}"
	return result, nil
}

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = StringArray{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		// Remove the curly braces and split by comma
		str := string(v)
		str = strings.Trim(str, "{}")
		if str == "" {
			*a = StringArray{}
			return nil
		}
		// Split and remove quotes
		parts := strings.Split(str, ",")
		result := make(StringArray, len(parts))
		for i, part := range parts {
			result[i] = strings.Trim(part, "\"")
		}
		*a = result
		return nil
	case string:
		return a.Scan([]byte(v))
	default:
		return errors.New("type assertion to []byte failed")
	}
}

type Chat struct {
	ID           string      `json:"id" gorm:"primaryKey"`
	ProjectID    string      `json:"project_id"`
	Participants StringArray `json:"participants" gorm:"type:text[]"`
	UserID       string      `json:"user_id"`
	IsGroupChat  bool        `json:"is_group_chat" gorm:"default:false"`
	GroupName    string      `json:"group_name"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

func (chat *Chat) BeforeCreate(tx *gorm.DB) (err error) {
	chat.ID = uuid.NewString()
	return
}

func (chat *Chat) ToChatResponse(users []UserResponse, lastMessage MessageResponse) ChatResponse {
	return ChatResponse{
		ID:           chat.ID,
		ProjectID:    chat.ProjectID,
		Users:        users,
		Participants: chat.Participants,
		IsGroupChat:  chat.IsGroupChat,
		GroupName:    chat.GroupName,
		LastMessage:  lastMessage,
		CreatedAt:    chat.CreatedAt,
		UpdatedAt:    chat.UpdatedAt,
	}
}

type ChatResponse struct {
	ID           string          `json:"id"`
	ProjectID    string          `json:"project_id"`
	Participants []string        `json:"participants"`
	Users        []UserResponse  `json:"users"`
	LastMessage  MessageResponse `json:"last_message"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	IsGroupChat  bool            `json:"is_group_chat"`
	GroupName    string          `json:"group_name"`
}

type CreateChatRequest struct {
	Participants []string `json:"participants" gorm:"type:text[]"`
	IsGroupChat  bool     `json:"is_group_chat" gorm:"default:false"`
	GroupName    string   `json:"group_name"`
}
