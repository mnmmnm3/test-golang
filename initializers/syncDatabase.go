package initializers

import "go-api/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
