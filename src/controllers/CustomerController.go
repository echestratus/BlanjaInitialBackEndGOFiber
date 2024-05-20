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

func GetAllCustomers(c *fiber.Ctx) error {
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
	customers := models.SelectAllCustomers(sort, keyword, limit, offset)
	totalData := models.CountDataCustomers()
	totalPage := math.Ceil(float64(totalData) / float64(limit))
	result := map[string]interface{}{
		"data":        customers,
		"currentPage": page,
		"limit":       limit,
		"totalData":   totalData,
		"totalPage":   totalPage,
	}

	return c.JSON(result)
}

func GetCustomerById(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	foundCustomer := models.SelectCustomerById(id)
	if foundCustomer.Name == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": fmt.Sprintf("Customer with ID %d NOT FOUND!", id),
		})
	}
	return c.JSON(foundCustomer)
}
func GetDetailCustomer(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Locals("ID").(string))
	if role := c.Locals("role").(string); role != "customer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	foundCustomer := models.SelectCustomerById(id)
	if foundCustomer.Name == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": fmt.Sprintf("Customer with ID %d NOT FOUND!", id),
		})
	}
	return c.JSON(foundCustomer)
}

func CreateNewCustomer(c *fiber.Ctx) error {
	var newCustomer models.Customer
	var customer map[string]interface{}
	if err := c.BodyParser(&customer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}
	customer = helpers.XSSMiddleware(customer)
	mapstructure.Decode(customer, &newCustomer)

	newCustomer.Role = "customer"

	if _, _, errCustomer, errSeller := models.FindEmail(newCustomer.Email); errCustomer == nil || errSeller == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email already taken",
		})
	}

	errors := helpers.ValidateStruct(newCustomer)
	if len(errors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(newCustomer.Password), bcrypt.DefaultCost)
	newCustomer.Password = string(hashPassword)

	err := models.PostCustomer(&newCustomer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to register Customer",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Register customer created successfully",
	})
}

func UpdateCustomerById(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var updatedCustomer models.APICustomer
	var customer map[string]interface{}
	if err := c.BodyParser(&customer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}
	customer = helpers.XSSMiddleware(customer)
	mapstructure.Decode(customer, &updatedCustomer)

	foundCustomer := models.SelectCustomerByIdAllField(id)
	if foundCustomer.Name == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": fmt.Sprintf("Customer with ID %d NOT FOUND!", id),
		})
	}
	if updatedCustomer.Name == "" {
		updatedCustomer.Name = foundCustomer.Name
	}
	if updatedCustomer.Email == "" {
		updatedCustomer.Email = foundCustomer.Email
	}
	if updatedCustomer.Role == "" {
		updatedCustomer.Role = foundCustomer.Role
	}
	// if updatedCustomer.Password == "" {
	// 	updatedCustomer.Password = foundCustomer.Password
	// }

	if tempCustomer, _, errCustomer, errSeller := models.FindEmail(updatedCustomer.Email); (errCustomer == nil && tempCustomer.ID != foundCustomer.ID) || (errSeller == nil) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email already taken",
		})
	}

	errors := helpers.ValidateStruct(updatedCustomer)
	if len(errors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	err := models.UpdateCustomer(id, &updatedCustomer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("Failed to update customer with ID %d", id),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"messsage": fmt.Sprintf("Customer with ID %d UPDATED successfully", id),
	})
}
func UpdateCustomer(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Locals("ID").(string))
	role, _ := c.Locals("role").(string)
	if role != "customer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	var updatedCustomer models.APICustomer
	var customer map[string]interface{}
	if err := c.BodyParser(&customer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}
	customer = helpers.XSSMiddleware(customer)
	mapstructure.Decode(customer, &updatedCustomer)

	foundCustomer := models.SelectCustomerByIdAllField(id)
	if foundCustomer.Name == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": fmt.Sprintf("Customer with ID %d NOT FOUND!", id),
		})
	}
	if updatedCustomer.Name == "" {
		updatedCustomer.Name = foundCustomer.Name
	}
	if updatedCustomer.Email == "" {
		updatedCustomer.Email = foundCustomer.Email
	}
	if updatedCustomer.Role == "" {
		updatedCustomer.Role = foundCustomer.Role
	}
	// if updatedCustomer.Password == "" {
	// 	updatedCustomer.Password = foundCustomer.Password
	// }

	if tempCustomer, _, errCustomer, errSeller := models.FindEmail(updatedCustomer.Email); (errCustomer == nil && tempCustomer.ID != foundCustomer.ID) || (errSeller == nil) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email already taken",
		})
	}

	errors := helpers.ValidateStruct(updatedCustomer)
	if len(errors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	err := models.UpdateCustomer(id, &updatedCustomer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("Failed to update customer with ID %d", id),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"messsage": fmt.Sprintf("Customer with ID %d UPDATED successfully", id),
	})
}

func DeleteCustomer(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	foundCustomer := models.SelectCustomerById(id)
	if foundCustomer.Name == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": fmt.Sprintf("Customer with ID %d NOT FOUND!", id),
		})
	}
	err := models.DeleteCustomer(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("Failed to DELETE Customer with ID %d", id),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Customer with ID %d DELETED successfully", id),
	})
}
