package models

import (
	"BlanjaInitialBackEndGOFiber/src/configs"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string   `json:"name" validate:"min=3,max=100"`
	Price       float64  `json:"price" validate:"gte=0"`
	Size        int      `json:"size" validate:"gte=0"`
	Condition   string   `json:"condition" validate:"min=3"`
	Description string   `json:"description"`
	Rating      int      `json:"rating"`
	Color       string   `json:"color"`
	Image       string   `json:"imageURL"`
	CategoryID  uint     `json:"category_id"`
	SellerID    uint     `json:"seller_id"`
	Category    Category `gorm:"foreignKey:CategoryID" validate:"omitempty"`
	Seller      Seller   `gorm:"foreignKey:SellerID" validate:"omitempty"`
}

func SelectAllProducts(sort string, name string, limit int, offset int) []*Product {
	var items []*Product
	name = "%" + name + "%"
	configs.DB.Preload("Category").Preload("Seller").Order(sort).Limit(limit).Offset(offset).Where("deleted_at IS NULL").Where("LOWER(name) LIKE LOWER(?)", name).Find(&items)
	return items
}
func SelectAllProductsBySellerId(seller_id int) []*Product {
	var items []*Product
	configs.DB.Preload("Category").Preload("Seller").Where("deleted_at IS NULL").Where("seller_id = ?", seller_id).Find(&items)
	return items
}
func CountDataProducts() int64 {
	var result int64
	configs.DB.Table("products").Where("deleted_at IS NULL").Count(&result)
	return result
}

func SelectProductById(id int) *Product {
	var item *Product
	configs.DB.Preload("Category").Preload("Seller").First(&item, "id = ?", id)
	return item
}

func PostProduct(item *Product) error {
	result := configs.DB.Create(&item)
	return result.Error
}

func UpdateProduct(id int, updatedProduct *APIProduct) error {
	result := configs.DB.Model(&Product{}).Where("id = ?", id).Updates(updatedProduct)
	return result.Error
}

func DeleteProduct(id int) error {
	result := configs.DB.Delete(&Product{}, "id = ?", id)
	return result.Error
}
