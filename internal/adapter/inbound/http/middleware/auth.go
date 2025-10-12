package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nitikhon/golang-inventory-system/internal/util"
)

// AuthMiddleware is a middleware that checks for a valid JWT in the Authorization header.
// If the JWT is valid, it extracts the user ID and sets it in the context.
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if the Authorization header is present
		authHeader := c.Get("Authorization")
		parts := strings.Split(authHeader, " ")

		// If the Authorization header is not present or does not contain a Bearer token, return 401
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Authorization header format"})
		}

		// Extract the token from the Authorization header
		tokenString := parts[1]

		// Parse the token using the secret key
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(util.GetEnvOrPanic("ACCESS_TOKEN_SECRET")), nil
		})

		// If the token is invalid or expired, return 401
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired JWT"})
		}

		// Extract the claims from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid JWT claims"})
		}

		// Check if the user ID is present in the claims
		if userIDFloat, ok := claims["user_id"].(float64); ok {
			c.Locals("user_id", uint(userIDFloat))
		}
		
		// Store JWT claims in the context for later use.
		c.Locals("user", claims)
		return c.Next()
	}
}

func AdminOrOwnerOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ok for type assertion
		claims, ok := c.Locals("user").(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "No user context found",
            })
		}
		
		isAdmin, ok := claims["is_admin"].(bool)
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": "Invalid data in the token",
            })
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": "Invalid data in the token",
            })
		}
		userID := uint(userIDFloat)

		targetID, err := c.ParamsInt("id")
        if err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": "Invalid user ID parameter",
            })
        }

		if !isAdmin && userID != uint(targetID) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": "Invalid data in the token",
            })
		}

		return c.Next();
	}
}

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ok for type assertion
		claims, ok := c.Locals("user").(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "No user context found",
            })
		}
		
		isAdmin, ok := claims["is_admin"].(bool)
		if !ok || !isAdmin {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Admin permission required",
            })
		}

		return c.Next();
	}
}