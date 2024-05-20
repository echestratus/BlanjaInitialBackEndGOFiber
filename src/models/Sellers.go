package models

import (
	"BlanjaInitialBackEndGOFiber/src/configs"

	"gorm.io/gorm"
)

type Seller struct {
	gorm.Model
	Name        string       `json:"name" validate:"min=3,max=50"`
	Email       string       `json:"email" validate:"min=12,max=100"`
	PhoneNumber string       `json:"phoneNumber" validate:"e164,min=10,max=15"`
	StoreName   string       `json:"storeName" validate:"min=3,max=50"`
	Password    string       `json:"password" validate:"min=5,max=100,passwordValidation"`
	Role        string       `json:"role" validate:"min=2,max=100"`
	Products    []APIProduct `json:"products"`
}
type APISeller struct {
	gorm.Model
	Name        string `json:"name" validate:"min=3,max=50"`
	Email       string `json:"email" validate:"min=12,max=100"`
	PhoneNumber string `json:"phoneNumber" validate:"e164,min=10,max=15"`
	StoreName   string `json:"storeName" validate:"min=3,max=50"`
	Role        string `json:"role" validate:"min=2,max=100"`
}

func SelectAllSellers(sort string, name string, limit int, offset int) []*Seller {
	var sellers []*Seller
	name = "%" + name + "%"
	configs.DB.Model(&Seller{}).Preload("Products", func(db *gorm.DB) *gorm.DB {
		var items []*APIProduct
		return db.Model(&Product{}).Find(&items)
	}).Order(sort).Limit(limit).Offset(offset).Omit("password").Where("deleted_at IS NULL").Where("LOWER(name) LIKE LOWER(?)", name).Find(&sellers)
	return sellers
}
func CountDataSellers() int64 {
	var result int64
	configs.DB.Table("sellers").Where("deleted_at IS NULL").Count(&result)
	return result
}

func SelectSellerById(id int) *Seller {
	var seller *Seller
	configs.DB.Model(&Seller{}).Preload("Products", func(db *gorm.DB) *gorm.DB {
		var items []*APIProduct
		return db.Model(&Product{}).Find(&items)
	}).Omit("password").First(&seller, "id = ?", id)
	return seller
}

func PostSeller(seller *Seller) error {
	result := configs.DB.Create(&seller)
	return result.Error
}

func UpdateSeller(id int, updatedSeller *APISeller) error {
	result := configs.DB.Model(&Seller{}).Where("id = ?", id).Updates(updatedSeller)
	return result.Error
}

func DeleteSeller(id int) error {
	result := configs.DB.Delete(&Seller{}, "id = ?", id)
	return result.Error
}
