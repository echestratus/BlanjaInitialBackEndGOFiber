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

		token, err := helpers.GenerateToken(os.Getenv("SECRET_KEY"), map[string]interface{}{"ID": ID, "email": customer.Email, "role": customer.Role})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate token",
			})
		}

		refreshToken, err := helpers.GenerateRefreshToken(os.Getenv("SECRET_KEY"), map[string]interface{}{"ID": ID, "email": customer.Email, "role": customer.Role})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate refresh Token",
			})
		}

		return c.JSON(fiber.Map{
			"message":      "Login Successfully",
			"email":        customer.Email,
			"token":        token,
			"refreshToken": refreshToken,
		})
	} else {
		if err := bcrypt.CompareHashAndPassword([]byte(seller.Password), []byte(input.Password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid Password",
			})
		}

		ID := strconv.Itoa(int(seller.ID))
		token, err := helpers.GenerateToken(os.Getenv("SECRET_KEY"), map[string]interface{}{"ID": ID, "email": seller.Email, "role": seller.Role})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate token",
			})
		}

		refreshToken, err := helpers.GenerateRefreshToken(os.Getenv("SECRET_KEY"), map[string]interface{}{"ID": ID, "email": seller.Email, "role": seller.Role})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate refresh Token",
			})
		}

		return c.JSON(fiber.Map{
			"message":      "Login Successfully",
			"email":        seller.Email,
			"token":        token,
			"refreshToken": refreshToken,
		})
	}
}

func RefreshToken(c *fiber.Ctx) error {
	var input struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	token, err := helpers.GenerateToken(os.Getenv("SECRET_KEY"), map[string]interface{}{"refreshToken": input.RefreshToken})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate access Token",
		})
	}

	refreshToken, err := helpers.GenerateRefreshToken(os.Getenv("SECRET_KEY"), map[string]interface{}{"refreshToken": input.RefreshToken})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate refresh Token",
		})
	}

	item := map[string]string{
		"token":        token,
		"refreshToken": refreshToken,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Refresh succeed",
		"data":    item,
	})
}
