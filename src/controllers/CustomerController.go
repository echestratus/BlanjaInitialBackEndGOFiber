package controllers

import (
	"BlanjaInitialBackEndGOFiber/src/models"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetAllCustomers(c *fiber.Ctx) error {
	customers := models.SelectAllCustomers()
	return c.JSON(customers)
}

func GetDetailCustomer(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
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
	if err := c.BodyParser(&newCustomer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}
	models.PostCustomer(&newCustomer)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "Customer created successfully",
		"customer": newCustomer,
	})
}

func UpdateCustomer(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var updatedCustomer models.Customer
	if err := c.BodyParser(&updatedCustomer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}
	foundCustomer := models.SelectCustomerById(id)
	if foundCustomer.Name == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": fmt.Sprintf("Customer with ID %d NOT FOUND!", id),
		})
	}
	err := models.UpdateCustomer(id, &updatedCustomer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("Failed to update customer with ID %d", id),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"messsage":        fmt.Sprintf("Customer with ID %d UPDATED successfully", id),
		"updatedCustomer": updatedCustomer,
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
