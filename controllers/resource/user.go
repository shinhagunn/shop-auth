package resource

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shinhagunn/shop-auth/config/collection"
	"github.com/shinhagunn/shop-auth/controllers"
	"github.com/shinhagunn/shop-auth/models"
	"github.com/shinhagunn/shop-auth/services"
)

const (
	wrongPasswordErr  = "Old password is incorrect"
	updatePasswordErr = "Update password failed"
)

type UpdatePassPayload struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func UpdatePassword(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)

	payload := new(UpdatePassPayload)

	if err := c.BodyParser(payload); err != nil {
		return c.JSON(controllers.FailedToParseBody)
	}

	if result := services.CheckPasswordHash(payload.OldPassword, user.Password); !result {
		return c.JSON(wrongPasswordErr)
	}

	user.Password = services.HashPassword(payload.NewPassword)

	if err := collection.User.Update(user); err != nil {
		return c.JSON(updatePasswordErr)
	}

	return c.Status(200).JSON(200)
}
