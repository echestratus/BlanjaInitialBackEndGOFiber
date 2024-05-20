package controllers

import (
	"BlanjaInitialBackEndGOFiber/src/helpers"
	"BlanjaInitialBackEndGOFiber/src/models"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
)

func GetAllCategories(c *fiber.Ctx) error {
	categories := models.SelectAllCategories()
	return c.JSON(categories)
}

func GetDetailCategory(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	category := models.SelectCategoryById(id)
	if category.Name == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": fmt.Sprintf("Category with id %d NOT FOUND!", id),
		})
	}
	return c.JSON(category)
}

func CreateCategory(c *fiber.Ctx) error {
	var category map[string]interface{}
	category = helpers.XSSMiddleware(category)
	var newCategory models.Category
	mapstructure.Decode(category, &newCategory)
	if err := c.BodyParser(&newCategory); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	errors := helpers.ValidateStruct(newCategory)
	if len(errors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	err := models.PostCategory(&newCategory)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to post category",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "Category created successfully",
		"category": newCategory,
	})
}

func UpdateCategory(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var updatedCategory models.Category
	var category map[string]interface{}
	category = helpers.XSSMiddleware(category)
	mapstructure.Decode(category, &updatedCategory)
	if err := c.BodyParser(&updatedCategory); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid body request",
		})
	}

	errors := helpers.ValidateStruct(updatedCategory)
	if len(errors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	foundCategory := models.SelectCategoryById(id)
	if foundCategory.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": fmt.Sprintf("Category with ID %d NOT FOUND!", id),
		})
	}

	err := models.UpdateCategory(id, &updatedCategory)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update category",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category updated successfully",
	})
}

func DeleteCategory(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	foundProduct := models.SelectCategoryById(id)
	if foundProduct.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": fmt.Sprintf("Category with ID %d NOT FOUND!", id),
		})
	}
	models.DeleteCategory(id)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Category with ID %d DELETED", id),
	})
}
