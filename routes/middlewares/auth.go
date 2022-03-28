package middlewares

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/shinhagunn/Shop-Watches/backend/config"
	"github.com/shinhagunn/Shop-Watches/backend/config/collection"
	"github.com/shinhagunn/Shop-Watches/backend/controllers"
	"github.com/shinhagunn/Shop-Watches/backend/models"
	"github.com/shinhagunn/Shop-Watches/backend/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func MustAuth(c *fiber.Ctx) error {
	path := strings.Split(c.Path(), "/api/v2/auth")[1]
	if strings.Contains(path, "/api/v2/myauth/identity") || strings.Contains(path, "/api/v2/myauth/public") {
		return c.SendStatus(200)
	}

	session, err := config.SessionStore.Get(c)

	if err != nil {
		return c.Status(500).JSON(controllers.FailedConnectToSessions)
	}

	uid := session.Get("uid")

	if uid == nil {
		return c.Status(401).JSON("Not logged in")
	}

	user := new(models.User)
	result := collection.User.FindOne(context.Background(), bson.M{"uid": uid})
	result.Decode(&user)

	if len(user.Email) == 0 {
		session.Destroy()
		session.Save()
	}

	return c.Next()
}

type Auth struct {
	UID         string   `json:"uid"`
	State       string   `json:"state"`
	Email       string   `json:"email"`
	Username    string   `json:"username"`
	Role        string   `json:"role"`
	ReferralUID string   `json:"referral_uid"`
	Level       int64    `json:"level"`
	Audience    []string `json:"aud,omitempty"`

	jwt.StandardClaims
}

func MustActive(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)

	if user.State != "Active" {
		return c.Status(422).JSON("User state must be active")
	}

	return c.Next()
}

func MustPending(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)

	if user.State != "Pending" {
		return c.Status(422).JSON("User state must be pending")
	}

	return c.Next()
}

func MustGuest(c *fiber.Ctx) error {
	session, err := config.SessionStore.Get(c)

	if err != nil {
		return c.Status(500).JSON(controllers.FailedConnectToSessions)
	}

	uid := session.Get("uid")

	if uid == nil {
		return c.Next()
	}

	return c.Status(422).JSON(controllers.MustBeGuest)
}

func CheckJWT(c *fiber.Ctx) error {
	jwt_auth, err := utils.CheckJWT(strings.Replace(c.Get("Authorization"), "Bearer ", "", -1))

	if err != nil {
		return c.Status(500).JSON(controllers.FailedToParseJWT)
	}

	user := new(models.User)

	result := collection.User.FindOne(context.Background(), bson.M{"uid": jwt_auth.UID})
	result.Decode(&user)

	if len(user.Email) == 0 {
		return c.Status(500).JSON(controllers.ServerInternalError)
	}

	c.Locals("CurrentUser", user)

	return c.Next()
}
