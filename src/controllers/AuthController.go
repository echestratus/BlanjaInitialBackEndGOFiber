package controllers

import (
	"BlanjaInitialBackEndGOFiber/src/helpers"
	"BlanjaInitialBackEndGOFiber/src/models"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func AuthLogin(c *fiber.Ctx) error {
	var input models.Login
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid body request",
		})
	}

	customer, seller, errCustomer, errSeller := models.FindEmail(input.Email)
	if errCustomer != nil && errSeller != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Email not found",
		})
	}
	if customer.Email != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(input.Password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid Password",
			})
		}

		ID := strconv.Itoa(int(customer.ID))
		token, err := helpers.GenerateToken(os.Getenv("SECRET_KEY"), ID, customer.Email, customer.Role)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate token",
			})
		}
		return c.JSON(fiber.Map{
			"message": "Login Successfully",
			"token":   token,
		})
	} else {
		if err := bcrypt.CompareHashAndPassword([]byte(seller.Password), []byte(input.Password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid Password",
			})
		}

		ID := strconv.Itoa(int(seller.ID))
		token, err := helpers.GenerateToken(os.Getenv("SECRET_KEY"), ID, seller.Email, seller.Role)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate token",
			})
		}
		return c.JSON(fiber.Map{
			"message": "Login Successfully",
			"token":   token,
		})
	}
}
