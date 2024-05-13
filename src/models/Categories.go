package models

import (
	"BlanjaInitialBackEndGOFiber/src/configs"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name     string       `json:"name"`
	Image    string       `json:"imageURL"`
	Products []APIProduct `json:"products"`
}

type APIProduct struct {
	gorm.Model
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Size        int     `json:"size"`
	Condition   string  `json:"condition"`
	Description string  `json:"description"`
	Rating      int     `json:"rating"`
	Color       string  `json:"color"`
	Image       string  `json:"imageURL"`
	CategoryID  uint    `json:"category_id"`
}

func SelectAllCategories() []*Category {
	var categories []*Category
	configs.DB.Preload("Products", func(db *gorm.DB) *gorm.DB {
		var items []*APIProduct
		return db.Model(&Product{}).Find(&items)
	}).Find(&categories)
	return categories
}

func SelectCategoryById(id int) *Category {
	var category *Category
	configs.DB.Preload("Products", func(db *gorm.DB) *gorm.DB {
		var items []*APIProduct
		return db.Model(&Product{}).Find(&items)
	}).First(&category, "id = ?", id)
	return category
}

func PostCategory(category *Category) error {
	result := configs.DB.Create(&category)
	return result.Error
}

func UpdateCategory(id int, updatedCategory *Category) error {
	result := configs.DB.Model(&Category{}).Where("id = ?", id).Updates(updatedCategory)
	return result.Error
}

func DeleteCategory(id int) error {
	result := configs.DB.Delete(&Category{}, "id = ?", id)
	return result.Error
}
