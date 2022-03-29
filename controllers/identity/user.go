package identity

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/shinhagunn/shop-auth/config"
	"github.com/shinhagunn/shop-auth/config/collection"
	"github.com/shinhagunn/shop-auth/controllers"
	"github.com/shinhagunn/shop-auth/models"
	"github.com/shinhagunn/shop-auth/services"
	"github.com/shinhagunn/shop-auth/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type RegisterPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *fiber.Ctx) error {
	payload := new(RegisterPayload)

	if err := c.BodyParser(payload); err != nil {
		return c.Status(422).JSON(controllers.FailedToParseBody)
	}

	user := &models.User{
		UID:      utils.RandomUID(),
		Email:    payload.Email,
		Password: services.HashPassword(payload.Password),
		State:    "Pending",
		Role:     "Member",
	}

	if err := collection.User.Create(user); err != nil {
		return c.Status(422).JSON(controllers.FailedConnectDataInDatabase)
	}

	services.EmailProducer("new-user", user)

	session, err := config.SessionStore.Get(c)

	if err != nil {
		return c.Status(500).JSON(controllers.ServerInternalError)
	}

	session.Set("uid", user.UID)
	session.Save()

	return c.SendStatus(200)
}

func ReSendEmailCode(c *fiber.Ctx) error {
	session, err := config.SessionStore.Get(c)

	if err != nil {
		return c.Status(500).JSON(controllers.FailedConnectToSessions)
	}

	uid := session.Get("uid")

	user := new(models.User)

	collection.User.FindOne(context.Background(), bson.M{"uid": uid}).Decode(user)

	if user.State != "Pending" {
		return c.Status(422).JSON("Must be Pending")
	}

	// Check old code and change it's status
	codes := []models.Code{}
	collection.Code.SimpleFind(&codes, bson.M{"user_id": user.ID})

	for _, code := range codes {
		code.State = "Delete"
		collection.Code.Update(&code)
	}

	// Create new code
	services.EmailProducer("new-user", user)

	return c.JSON(200)
}

type VerifyPayload struct {
	Code string `json:"code"`
}

func VerificationCode(c *fiber.Ctx) error {
	payload := new(VerifyPayload)

	if err := c.BodyParser(payload); err != nil {
		return c.Status(422).JSON(controllers.FailedToParseBody)
	}

	session, err := config.SessionStore.Get(c)

	if err != nil {
		return c.Status(500).JSON(controllers.FailedConnectToSessions)
	}

	uid := session.Get("uid")

	user := new(models.User)

	collection.User.FindOne(context.Background(), bson.M{"uid": uid}).Decode(user)

	if user.State != "Pending" {
		return c.Status(422).JSON("Must be Pending")
	}

	code := new(models.Code)
	collection.Code.FindOne(context.Background(), bson.M{"user_id": user.ID, "state": "Active"}).Decode(&code)

	if code.Code == "" {
		return c.Status(422).JSON(controllers.FailedConnectDataInDatabase)
	}

	if code.Code != payload.Code {
		return c.Status(422).JSON(controllers.FailedConnectDataInDatabase)
	}

	code.State = "Delete"
	if err := collection.Code.Update(code); err != nil {
		return c.Status(422).JSON(controllers.FailedConnectDataInDatabase)
	}

	user.State = "Active"
	if err := collection.User.Update(user); err != nil {
		return c.Status(422).JSON(controllers.FailedConnectDataInDatabase)
	}

	if err != nil {
		return c.Status(500).JSON(controllers.FailedConnectToSessions)
	}

	session.Set("uid", user.UID)
	session.Save()

	return c.JSON(200)
}
