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
	"golang.org/x/crypto/bcrypt"
)

func GetAllSellers(c *fiber.Ctx) error {
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
	sellers := models.SelectAllSellers(sort, keyword, limit, offset)
	totalData := models.CountDataSellers()
	totalPage := math.Ceil(float64(totalData) / float64(limit))
	result := map[string]interface{}{
		"data":        sellers,
		"currentPage": page,
		"limit":       limit,
		"totalData":   totalData,
		"totalPage":   totalPage,
	}

	return c.JSON(result)
}
func GetSellerById(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	foundSeller := models.SelectSellerById(id)
	if foundSeller.Email == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": fmt.Sprintf("Seller with ID %d NOT FOUND!", id),
		})
	}
	return c.JSON(foundSeller)
}
func GetDetailSeller(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Locals("ID").(string))
	if role := c.Locals("role").(string); role != "seller" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	foundSeller := models.SelectSellerById(id)
	if foundSeller.Email == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": fmt.Sprintf("Seller with ID %d NOT FOUND!", id),
		})
	}
	return c.JSON(foundSeller)
}

func CreateNewSeller(c *fiber.Ctx) error {
	var newSeller models.Seller
	var seller map[string]interface{}
	if err := c.BodyParser(&seller); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}
	seller = helpers.XSSMiddleware(seller)
	mapstructure.Decode(seller, &newSeller)

	newSeller.Role = "seller"

	if _, _, errCustomer, errSeller := models.FindEmail(newSeller.Email); errCustomer == nil || errSeller == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email already taken",
		})
	}

	errors := helpers.ValidateStruct(newSeller)
	if len(errors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(newSeller.Password), bcrypt.DefaultCost)
	newSeller.Password = string(hashPassword)

	err := models.PostSeller(&newSeller)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to register Seller",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Register seller created successfully",
	})
}

func UpdateSellerById(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var updatedSeller models.APISeller
	var seller map[string]interface{}
	if err := c.BodyParser(&seller); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}
	seller = helpers.XSSMiddleware(seller)
	mapstructure.Decode(seller, &updatedSeller)

	foundSeller := models.SelectSellerById(id)
	if foundSeller.Email == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": fmt.Sprintf("Seller with ID %d NOT FOUND!", id),
		})
	}
	if updatedSeller.Name == "" {
		updatedSeller.Name = foundSeller.Name
	}
	if updatedSeller.Email == "" {
		updatedSeller.Email = foundSeller.Email
	}
	if updatedSeller.PhoneNumber == "" {
		updatedSeller.PhoneNumber = foundSeller.PhoneNumber
	}
	if updatedSeller.StoreName == "" {
		updatedSeller.StoreName = foundSeller.StoreName
	}
	if updatedSeller.Role == "" {
		updatedSeller.Role = foundSeller.Role
	}

	if _, tempSeller, errCustomer, errSeller := models.FindEmail(updatedSeller.Email); (errSeller == nil && tempSeller.ID != foundSeller.ID) || (errCustomer == nil) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email already taken",
		})
	}

	errors := helpers.ValidateStruct(updatedSeller)
	if len(errors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	err := models.UpdateSeller(id, &updatedSeller)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("Failed to update seller with ID %d", id),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"messsage": fmt.Sprintf("Seller with ID %d UPDATED successfully", id),
	})
}
func UpdateSeller(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Locals("ID").(string))
	role, _ := c.Locals("role").(string)
	if role != "seller" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	var updatedSeller models.APISeller
	var seller map[string]interface{}
	if err := c.BodyParser(&seller); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}
	seller = helpers.XSSMiddleware(seller)
	mapstructure.Decode(seller, &updatedSeller)

	foundSeller := models.SelectSellerById(id)
	if foundSeller.Email == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": fmt.Sprintf("Seller with ID %d NOT FOUND!", id),
		})
	}
	if updatedSeller.Name == "" {
		updatedSeller.Name = foundSeller.Name
	}
	if updatedSeller.Email == "" {
		updatedSeller.Email = foundSeller.Email
	}
	if updatedSeller.PhoneNumber == "" {
		updatedSeller.PhoneNumber = foundSeller.PhoneNumber
	}
	if updatedSeller.StoreName == "" {
		updatedSeller.StoreName = foundSeller.StoreName
	}
	if updatedSeller.Role == "" {
		updatedSeller.Role = foundSeller.Role
	}

	if _, tempSeller, errCustomer, errSeller := models.FindEmail(updatedSeller.Email); (errSeller == nil && tempSeller.ID != foundSeller.ID) || (errCustomer == nil) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email already taken",
		})
	}

	errors := helpers.ValidateStruct(updatedSeller)
	if len(errors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	err := models.UpdateSeller(id, &updatedSeller)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("Failed to update seller with ID %d", id),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"messsage": fmt.Sprintf("Seller with ID %d UPDATED successfully", id),
	})
}

func DeleteSeller(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	foundSeller := models.SelectSellerById(id)
	if foundSeller.Email == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": fmt.Sprintf("Seller with ID %d NOT FOUND!", id),
		})
	}
	foundProducts := models.SelectAllProductsBySellerId(id)
	for _, foundProduct := range foundProducts {
		err := models.DeleteProduct(int(foundProduct.ID))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": fmt.Sprintf("Failed to delete product with ID %d", id),
			})
		}
	}

	err := models.DeleteSeller(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("Failed to DELETE Seller with ID %d", id),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Seller with ID %d DELETED successfully", id),
	})
}
