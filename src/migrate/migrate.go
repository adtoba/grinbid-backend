package migrate

import (
	"log"

	"github.com/adtoba/grinbid-backend/src/models"
	"gorm.io/gorm"
)

func Migrate(DB *gorm.DB) {
	DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Listing{},
		&models.Wallet{},
		&models.Transaction{},
		&models.Chat{},
		&models.Message{},
		&models.MessageSeen{},
		&models.Session{},
	)

	log.Println("Database migrated successfully")
}
