package models

import (
	"BlanjaInitialBackEndGOFiber/src/configs"

	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email,min=12,max=100"`
	Password string `json:"password" validate:"required,min=5,max=100,passwordValidation"`
	Role     string `json:"role" validate:"required,min=2,max=100"`
}

type APICustomer struct {
	gorm.Model
	Name  string `json:"name" validate:"required,min=3,max=50"`
	Email string `json:"email" validate:"required,email,min=12,max=100"`
	Role  string `json:"role" validate:"required,min=2,max=100"`
}

func SelectAllCustomers(sort string, name string, limit int, offset int) []*APICustomer {
	var items []*APICustomer
	name = "%" + name + "%"
	configs.DB.Model(&Customer{}).Order(sort).Limit(limit).Offset(offset).Where("deleted_at IS NULL").Where("LOWER(name) LIKE LOWER(?)", name).Find(&items)
	return items
}
func CountDataCustomers() int64 {
	var result int64
	configs.DB.Table("customers").Where("deleted_at IS NULL").Count(&result)
	return result
}

func SelectCustomerById(id int) *APICustomer {
	var item *APICustomer
	configs.DB.Model(&Customer{}).First(&item, "id = ?", id)
	return item
}
func SelectCustomerByIdAllField(id int) *Customer {
	var item *Customer
	configs.DB.Model(&Customer{}).First(&item, "id = ?", id)
	return item
}

func PostCustomer(item *Customer) error {
	result := configs.DB.Create(&item)
	return result.Error
}

func UpdateCustomer(id int, updatedCustomer *APICustomer) error {
	result := configs.DB.Model(&Customer{}).Where("id = ?", id).Updates(updatedCustomer)
	return result.Error
}

func DeleteCustomer(id int) error {
	result := configs.DB.Delete(&Customer{}, "id = ?", id)
	return result.Error
}
