package middleware

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func Auth(c *fiber.Ctx) error {
	token := c.Cookies("jwt")
	if token == "" {
		authHeader := c.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized - No token provided",
		})
	}

	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY_JWT")), nil
	})
	if err != nil || !parsedToken.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized - Invalid token",
		})
	}
	log.Println(claims)
	username, ok := claims["username"].(string)
	if !ok || username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request - Username not found or invalid type",
		})
	}

	role, ok := claims["role"].(string)
	if !ok || role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request - Role not found or invalid type",
		})
	}

	userIDValue, ok := claims["user_id"]
	if !ok || userIDValue == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request - User ID not found or invalid type",
		})
	}

	var userID string
	switch v := userIDValue.(type) {
	case string:
		userID = v
	case float64:
		userID = fmt.Sprintf("%.0f", v)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request - User ID has invalid type",
		})
	}

	parkNo, ok := claims["parkno"].(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request - Park number not found in token",
		})
	}

	c.Locals("parkno", parkNo)
	c.Locals("user_id", userID)
	c.Locals("username", username)
	c.Locals("role", role)

	return c.Next()
}
