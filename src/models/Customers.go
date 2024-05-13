package models

import (
	"BlanjaInitialBackEndGOFiber/src/configs"

	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type APICustomer struct {
	gorm.Model
	Name  string `json:"name"`
	Email string `json:"email"`
}

func SelectAllCustomers() []*APICustomer {
	var items []*APICustomer
	configs.DB.Model(&Customer{}).Find(&items)
	return items
}

func SelectCustomerById(id int) *APICustomer {
	var item *APICustomer
	configs.DB.Model(&Customer{}).First(&item, "id = ?", id)
	return item
}

func PostCustomer(item *Customer) error {
	result := configs.DB.Create(&item)
	return result.Error
}

func UpdateCustomer(id int, updatedCustomer *Customer) error {
	result := configs.DB.Model(&Customer{}).Where("id = ?", id).Updates(updatedCustomer)
	return result.Error
}

func DeleteCustomer(id int) error {
	result := configs.DB.Delete(&Customer{}, "id = ?", id)
	return result.Error
}
