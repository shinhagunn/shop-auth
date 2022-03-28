package resource

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shinhagunn/Shop-Watches/backend/config/collection"
	"github.com/shinhagunn/Shop-Watches/backend/controllers"
	"github.com/shinhagunn/Shop-Watches/backend/models"
	"github.com/shinhagunn/Shop-Watches/backend/services"
)

const (
	wrongPasswordErr  = "Old password is incorrect"
	updatePasswordErr = "Update password failed"
)

func GetUserCurrent(c *fiber.Ctx) error {

	user := c.Locals("CurrentUser").(*models.User)

	return c.JSON(user)
}

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

type ProfilePayload struct {
	Fullname string `json:"fullname"`
	Age      string `json:"age"`
	Address  string `json:"address"`
	Gender   string `json:"gender"`
	Phone    string `json:"phone"`
}

func UpdateUserProfile(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)
	payload := new(ProfilePayload)

	if err := c.BodyParser(payload); err != nil {
		return c.Status(422).JSON(controllers.FailedToParseBody)
	}

	if payload.Fullname == "" {
		user.UserProfile.Fullname = "none"
	} else {
		user.UserProfile.Fullname = payload.Fullname
	}

	if payload.Address == "" {
		user.UserProfile.Address = "none"
	} else {
		user.UserProfile.Address = payload.Address
	}

	if payload.Age == "" {
		user.UserProfile.Age = "none"
	} else {
		user.UserProfile.Age = payload.Age
	}

	if payload.Gender == "" {
		user.UserProfile.Gender = "none"
	} else {
		user.UserProfile.Gender = payload.Gender
	}

	if payload.Phone == "" {
		user.UserProfile.Phone = "none"
	} else {
		user.UserProfile.Phone = payload.Phone
	}

	if err := collection.User.Update(user); err != nil {
		return c.Status(422).JSON(controllers.FailedConnectDataInDatabase)
	}

	return c.JSON(200)
}
