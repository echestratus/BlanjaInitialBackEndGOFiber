package controllers

import (
	"BlanjaInitialBackEndGOFiber/src/helpers"
	"BlanjaInitialBackEndGOFiber/src/models"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
)

func GetAllProducts(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 5
	}
	offset := (page - 1) * limit
	sort := c.Query("sort")
	if sort == "" {
		sort = "ASC"
	}
	sortBy := c.Query("sortBy")
	if sortBy == "" {
		sortBy = "name"
	}
	sort = sortBy + " " + strings.ToUpper(sort)
	keyword := c.Query("search")
	products := models.SelectAllProducts(sort, keyword, limit, offset)
	totalData := models.CountDataProducts()
	totalPage := math.Ceil(float64(totalData) / float64(limit))
	result := map[string]interface{}{
		"data":        products,
		"currentPage": page,
		"limit":       limit,
		"totalData":   totalData,
		"totalPage":   totalPage,
	}

	return c.JSON(result)
}

func GetDetailProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	foundProduct := models.SelectProductById(id)
	if foundProduct.Name == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": fmt.Sprintf("Product with ID %d NOT FOUND!", id),
		})
	}
	return c.JSON(foundProduct)
}

func CreateNewProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Locals("ID").(string))
	if role := c.Locals("role").(string); role != "seller" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	var newProduct models.Product
	var product map[string]interface{}
	product = helpers.XSSMiddleware(product)
	mapstructure.Decode(product, &newProduct)
	if err := c.BodyParser(&newProduct); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid body request",
		})
	}

	newProduct.SellerID = uint(id)

	errors := helpers.ValidateStruct(newProduct)
	if len(errors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	err := models.PostProduct(&newProduct)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create new product",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Product created successfully",
		"product": newProduct,
	})
}

func UpdateProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var updatedProduct models.APIProduct
	var product map[string]interface{}
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}
	product = helpers.XSSMiddleware(product)
	mapstructure.Decode(product, &updatedProduct)

	errors := helpers.ValidateStruct(updatedProduct)
	if len(errors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	foundProduct := models.SelectProductById(id)
	if foundProduct.Name == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": fmt.Sprintf("Product with ID %d NOT FOUND!", id),
		})
	}
	err := models.UpdateProduct(id, &updatedProduct)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("Failed to update product with ID %d", id),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product updated successfully",
	})
}

func DeleteProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	foundProduct := models.SelectProductById(id)
	if foundProduct.Name == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": fmt.Sprintf("Product with ID %d NOT FOUND!", id),
		})
	}
	err := models.DeleteProduct(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("Failed to delete product with ID %d", id),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Product with ID %d DELETED", id),
	})
}
