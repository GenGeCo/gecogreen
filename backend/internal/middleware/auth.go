package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/gecogreen/backend/internal/auth"
	"github.com/gecogreen/backend/internal/models"
	"github.com/gecogreen/backend/internal/repository"
)

func AuthMiddleware(jwtManager *auth.JWTManager, userRepo *repository.UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token mancante"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Formato token non valido"})
		}

		claims, err := jwtManager.ValidateToken(parts[1])
		if err != nil {
			if err == auth.ErrExpiredToken {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token scaduto"})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token non valido"})
		}

		if claims.TokenType != auth.AccessToken {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Tipo token non valido"})
		}

		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		user, err := userRepo.GetByID(ctx, claims.UserID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Utente non trovato"})
		}

		if !user.IsActive() {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Account non attivo"})
		}

		c.Locals("user", user)
		c.Locals("userID", user.ID)
		return c.Next()
	}
}

func RoleMiddleware(allowedRoles ...models.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*models.User)
		for _, role := range allowedRoles {
			if user.HasRole(role) {
				return c.Next()
			}
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Non autorizzato"})
	}
}

func SellerOnly() fiber.Handler {
	return RoleMiddleware(models.RoleSeller, models.RoleAdmin)
}

func AdminOnly() fiber.Handler {
	return RoleMiddleware(models.RoleAdmin)
}
