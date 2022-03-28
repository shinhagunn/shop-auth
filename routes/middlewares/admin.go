package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shinhagunn/Shop-Watches/backend/models"
)

func MustAdmin(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)

	if user.Role != "Admin" {
		return c.JSON("Not user Admin")
	}

	return c.Next()
}
