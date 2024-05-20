package helpers

import (
	"BlanjaInitialBackEndGOFiber/src/configs"
	"BlanjaInitialBackEndGOFiber/src/models"
)

func Migration() {
	configs.DB.AutoMigrate(&models.Product{}, &models.Category{}, &models.Customer{}, &models.Seller{})
}
