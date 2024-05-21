package middlewares

import (
	"fmt"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func ExtractToken(c *fiber.Ctx) string {
	bearerToken := c.Get("Authorization")
	if strings.HasPrefix(bearerToken, "Bearer ") {
		return strings.TrimPrefix(bearerToken, "Bearer ")
	}
	return ""
}

func JwtMiddleware() fiber.Handler {
	secretKey := os.Getenv("SECRET_KEY")
	return func(c *fiber.Ctx) error {
		tokenString := ExtractToken(c)
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
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
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		payload := token.Claims.(jwt.MapClaims)
		// id := payload["ID"].(string)
		// email := payload["email"].(string)
		// role := payload["role"].(string)
		// c.Locals("email", email)
		// c.Locals("ID", id)
		// c.Locals("role", role)
		for key, value := range payload {
			c.Locals(key, value)
		}

		return c.Next()
	}
}
