package migrate

import (
	"log"

	"github.com/adtoba/grinbid-backend/src/models"
	"gorm.io/gorm"
)

func Migrate(DB *gorm.DB) {
	DB.AutoMigrate(
		&models.Category{},
		&models.Listing{},
		&models.Wallet{},
		&models.Transaction{},
	)

	log.Println("Database migrated successfully")
}
