package middleware

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func OnlyAdmin(db *gorm.DB, fn fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("userRole").(string)
		if !ok || role != "Admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
		}
		return fn(c)
	}
}

func OnlyModerator(db *gorm.DB, fn fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("userRole").(string)
		if !ok || (role != "Moderator" && role != "Admin") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
		}
		return fn(c)
	}
}

func OnlyUser(db *gorm.DB, fn fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("userRole").(string)
		if !ok || role != "User" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
		}
		return fn(c)
	}
}
