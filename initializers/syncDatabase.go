package initializers

import "github.com/calleros/go-jwt/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
