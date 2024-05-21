package helpers

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func GenerateToken(secretKey string, payload map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	for key, value := range payload {
		if key == "refreshToken" {
			var c *fiber.Ctx
			tokenString := value.(string)
			if tokenString == "" {
				return "", c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Unauthorized",
				})
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secretKey), nil
			})

			if err != nil {
				return "", c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Unauthorized",
				})
			}

			payload := token.Claims.(jwt.MapClaims)
			for key, value := range payload {
				claims[key] = value
			}
		} else {
			claims[key] = value
		}

	}

	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	return token.SignedString([]byte(secretKey))
}

func GenerateRefreshToken(secretKey string, payload map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	for key, value := range payload {
		if key == "refreshToken" {
			var c *fiber.Ctx
			tokenString := value.(string)
			if tokenString == "" {
				return "", c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Unauthorized",
				})
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secretKey), nil
			})

			if err != nil {
				return "", c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Unauthorized",
				})
			}

			payload := token.Claims.(jwt.MapClaims)
			for key, value := range payload {
				claims[key] = value
			}
		} else {
			claims[key] = value
		}

	}

	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	return token.SignedString([]byte(secretKey))
}
