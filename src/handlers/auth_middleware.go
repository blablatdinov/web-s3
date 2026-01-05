package handlers

import (
	"strings"

	"github.com/blablatdinov/web-s3/src/srv"
	fiber "github.com/gofiber/fiber/v2"
)

const (
	UserIDKey   = "user_id"
	UsernameKey = "username"
)

func AuthMiddleware(userAuthSrv srv.UserAuth) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format. Expected: Bearer <token>",
			})
		}

		token := parts[1]
		claims, err := userAuthSrv.ExtractClaims(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}
		if userID, ok := claims["user_id"].(float64); ok {
			c.Locals(UserIDKey, int(userID))
		}
		if username, ok := claims["username"].(string); ok {
			c.Locals(UsernameKey, username)
		}
		return c.Next()
	}
}

func GetUserID(c *fiber.Ctx) (int, bool) {
	userID, ok := c.Locals(UserIDKey).(int)
	return userID, ok
}

func GetUsername(c *fiber.Ctx) (string, bool) {
	username, ok := c.Locals(UsernameKey).(string)
	return username, ok
}

