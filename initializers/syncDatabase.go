package initializers

import "github.com/calleros/sich/models"

func SyncDatabase() {
	DB.AutoMigrate(
		&models.User{},
		&models.Profile{},
	)
}
