package public

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shinhagunn/Shop-Watches/backend/config/collection"
	"github.com/shinhagunn/Shop-Watches/backend/controllers"
	"github.com/shinhagunn/Shop-Watches/backend/models"
)

type MessagePayload struct {
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Message  string `json:"message"`
}

func CustommerMessage(c *fiber.Ctx) error {
	payload := new(MessagePayload)

	if err := c.BodyParser(payload); err != nil {
		return c.Status(422).JSON(controllers.FailedToParseBody)
	}

	custommer := &models.Custommer{
		Fullname: payload.Fullname,
		Email:    payload.Email,
		Phone:    payload.Phone,
		Address:  payload.Address,
		Message:  payload.Message,
	}

	if err := collection.Custommer.Create(custommer); err != nil {
		return c.Status(500).JSON(controllers.ServerInternalError)
	}

	return c.JSON(custommer)
}
