package identity

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/shinhagunn/Shop-Watches/backend/config"
	"github.com/shinhagunn/Shop-Watches/backend/config/collection"
	"github.com/shinhagunn/Shop-Watches/backend/controllers"
	"github.com/shinhagunn/Shop-Watches/backend/models"
	"github.com/shinhagunn/Shop-Watches/backend/services"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	passwordWrongErr = "password is incorrect"
)

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *fiber.Ctx) error {
	payload := new(LoginPayload)

	if err := c.BodyParser(payload); err != nil {
		return c.Status(422).JSON(controllers.FailedToParseBody)
	}

	user := new(models.User)

	collection.User.FindOne(context.Background(), bson.M{"email": payload.Email}).Decode(&user)

	if user.Role == "" {
		return c.Status(422).JSON(controllers.FailedConnectDataInDatabase)
	}

	if user.State == "Delete" {
		return c.Status(422).JSON(controllers.FailedConnectDataInDatabase)
	}

	if user.State == "Banned" {
		return c.Status(422).JSON(controllers.FailedConnectDataInDatabase)
	}

	if result := services.CheckPasswordHash(payload.Password, user.Password); !result {
		return c.Status(422).JSON(passwordWrongErr)
	}

	session, err := config.SessionStore.Get(c)

	if err != nil {
		return c.Status(500).JSON(controllers.ServerInternalError)
	}

	session.Set("uid", user.UID)
	session.Save()

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	session, err := config.SessionStore.Get(c)

	if err != nil {
		return c.Status(500).JSON(controllers.ServerInternalError)
	}

	session.Destroy()
	session.Save()

	return c.JSON(200)
}
