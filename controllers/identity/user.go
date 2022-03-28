package identity

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/shinhagunn/Shop-Watches/backend/config"
	"github.com/shinhagunn/Shop-Watches/backend/config/collection"
	"github.com/shinhagunn/Shop-Watches/backend/controllers"
	"github.com/shinhagunn/Shop-Watches/backend/models"
	"github.com/shinhagunn/Shop-Watches/backend/services"
	"github.com/shinhagunn/Shop-Watches/backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		UID:         utils.RandomUID(),
		Email:       payload.Email,
		Password:    services.HashPassword(payload.Password),
		ChatIDs:     []primitive.ObjectID{},
		State:       "Pending",
		Role:        "Member",
		UserProfile: models.UserProfile{},
	}

	if err := user.Collection().Create(user); err != nil {
		return c.Status(422).JSON(controllers.FailedConnectDataInDatabase)
	}

	services.EmailProducer("new-user", user)

	session, err := config.SessionStore.Get(c)

	if err != nil {
		return c.Status(500).JSON(controllers.FailedConnectToSessions)
	}

	session.Set("uid", user.UID)
	session.Save()

	return c.JSON(user)
}

type EmailPayload struct {
	Email string `json:"email"`
}

func ReSendEmailCode(c *fiber.Ctx) error {
	payload := new(EmailPayload)

	if err := c.BodyParser(payload); err != nil {
		return c.Status(422).JSON(controllers.FailedToParseBody)
	}

	user := new(models.User)
	collection.User.FindOne(context.Background(), bson.M{"email": payload.Email}).Decode(user)

	if user.Email == "" {
		return c.Status(422).JSON(controllers.FailedConnectDataInDatabase)
	}

	// Check old code and change it's status
	codes := []models.Code{}
	collection.Code.SimpleFind(&codes, bson.M{"user_id": user.ID})

	for _, code := range codes {
		code.State = "Delete"
		collection.Code.Update(&code)
	}

	// Create new code
	randomCode := utils.RandomCode()
	services.EmailProducer("new-user", user)

	code := &models.Code{
		UserID:         user.ID,
		Code:           randomCode,
		CodeExpiration: time.Now().Add(5 * time.Minute),
		State:          "Active",
	}

	if err := collection.Code.Create(code); err != nil {
		return c.Status(422).JSON(controllers.FailedConnectDataInDatabase)
	}

	return c.JSON(200)
}

type VerifyPayload struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func VerificationCode(c *fiber.Ctx) error {
	payload := new(VerifyPayload)

	if err := c.BodyParser(payload); err != nil {
		return c.Status(422).JSON(controllers.FailedToParseBody)
	}

	user := new(models.User)
	collection.User.FindOne(context.Background(), bson.M{"email": payload.Email}).Decode(&user)
	if user.Email == "" {
		return c.Status(422).JSON(controllers.FailedConnectDataInDatabase)
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

	session, err := config.SessionStore.Get(c)

	if err != nil {
		return c.Status(500).JSON(controllers.FailedConnectToSessions)
	}

	session.Set("uid", user.UID)
	session.Save()

	return c.JSON(200)
}
