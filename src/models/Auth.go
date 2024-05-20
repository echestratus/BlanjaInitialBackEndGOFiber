package models

import (
	"BlanjaInitialBackEndGOFiber/src/configs"

	"gorm.io/gorm"
)

type Login struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
}

func FindEmail(email string) (*Customer, *Seller, error, error) {
	var customer *Customer
	var seller *Seller
	resultCustomer := configs.DB.Where("email = ?", email).First(&customer)
	resultSeller := configs.DB.Where("email = ?", email).First(&seller)
	return customer, seller, resultCustomer.Error, resultSeller.Error
}
