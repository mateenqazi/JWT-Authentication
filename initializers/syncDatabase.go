package initializers

import "github.com/mateenqazi/jwt-authenication/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
