package middlewares

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/shinhagunn/shop-auth/config"
	"github.com/shinhagunn/shop-auth/config/collection"
	"github.com/shinhagunn/shop-auth/controllers"
	"github.com/shinhagunn/shop-auth/models"
	"github.com/shinhagunn/shop-auth/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func MustAuth(c *fiber.Ctx) error {
	path := strings.Replace(c.Path(), "/api/v2/auth", "", 1)
	if strings.Contains(path, "/api/v2/myauth/identity") || strings.Contains(path, "/api/v2/myauth/public") || strings.Contains(path, "api/v2/product/public") {
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

	return c.Status(422).JSON("Must be guest")
}

func CheckRequest(c *fiber.Ctx) error {
	jwt_auth, err := utils.CheckJWT(strings.Replace(c.Get("Authorization"), "Bearer ", "", -1))

	if err != nil {
		return c.Status(500).JSON(controllers.FailedToParseJWT)
	}

	user := new(models.User)

	if err := collection.User.FindOne(context.Background(), bson.M{"uid": jwt_auth.UID}).Decode(&user); err != nil {
		return c.Status(500).JSON(controllers.ServerInternalError)
	}

	if len(user.Email) == 0 {
		return c.Status(500).JSON(controllers.ServerInternalError)
	}

	c.Locals("CurrentUser", user)

	return c.Next()
}
