package initializers

import "github.com/mohamedfawas/user-management-system/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{}, &models.Admin{})
}
