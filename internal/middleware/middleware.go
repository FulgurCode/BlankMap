package middleware

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func SetupCORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Accept", "Authorization", "Content-Type"},
	})
}

func RequireAuth(jwtSecret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing or invalid authorization header",
			})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token claims",
			})
		}

		userID, err := uuid.Parse(claims["sub"].(string))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid user id in token",
			})
		}

		c.Locals("userID", userID)
		return c.Next()
	}
}
func SetupLogger() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		// Log incoming request
		log.Printf("\n→ %s %s\n  Headers: %v\n  Body: %s",
			c.Method(),
			c.OriginalURL(),
			c.GetReqHeaders()["Authorization"],
			string(c.Body()),
		)

		// Process request
		err := c.Next()

		// Log outgoing response
		log.Printf("← %s %s | %d | %s\n  Body: %s",
			c.Method(),
			c.OriginalURL(),
			c.Response().StatusCode(),
			time.Since(start),
			string(c.Response().Body()),
		)

		return err
	}
}
