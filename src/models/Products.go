package models

import (
	"BlanjaInitialBackEndGOFiber/src/configs"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string   `json:"name"`
	Price       float64  `json:"price"`
	Size        int      `json:"size"`
	Condition   string   `json:"condition"`
	Description string   `json:"description"`
	Rating      int      `json:"rating"`
	Color       string   `json:"color"`
	Image       string   `json:"imageURL"`
	CategoryID  uint     `json:"category_id"`
	Category    Category `gorm:"foreignKey:CategoryID"`
}

func SelectAllProducts() []*Product {
	var items []*Product
	configs.DB.Preload("Category").Find(&items)
	return items
}

func SelectProductById(id int) *Product {
	var item *Product
	configs.DB.Preload("Category").First(&item, "id = ?", id)
	return item
}

func PostProduct(item *Product) error {
	result := configs.DB.Create(&item)
	return result.Error
}

func UpdateProduct(id int, updatedProduct *Product) error {
	result := configs.DB.Model(&Product{}).Where("id = ?", id).Updates(updatedProduct)
	return result.Error
}

func DeleteProduct(id int) error {
	result := configs.DB.Delete(&Product{}, "id = ?", id)
	return result.Error
}
