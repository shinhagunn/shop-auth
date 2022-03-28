package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/shinhagunn/shop-auth/config"
	"github.com/shinhagunn/shop-auth/config/collection"
	"github.com/shinhagunn/shop-auth/controllers"
	"github.com/shinhagunn/shop-auth/controllers/identity"
	"github.com/shinhagunn/shop-auth/controllers/resource"
	"github.com/shinhagunn/shop-auth/models"
	"github.com/shinhagunn/shop-auth/routes/middlewares"
	"github.com/shinhagunn/shop-auth/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func InitRouter() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	// /api/v2/auth/*
	app.All("/api/v2/auth/*", middlewares.MustAuth, func(c *fiber.Ctx) error {
		session, err := config.SessionStore.Get(c)

		if err != nil {
			return c.Status(500).JSON(controllers.ServerInternalError)
		}

		uid := session.Get("uid")

		if uid == nil {
			return c.Status(500).JSON(controllers.ServerInternalError)
		}

		user := new(models.User)

		collection.User.FindOne(context.Background(), bson.M{"uid": uid}).Decode(user)

		if len(user.Email) == 0 {
			return c.Status(500).JSON(controllers.ServerInternalError)
		}

		if user.State != "Active" {
			return c.Status(401).JSON(controllers.AuthZPermissionDeniedErr)
		}

		jwt_token, err := utils.GenerateJWT(user)

		if err != nil {
			return c.Status(500).JSON(controllers.ServerInternalError)
		}

		jwt_token = "Bearer " + jwt_token
		c.Set("Authorization", jwt_token)

		return c.SendStatus(200)
	})

	api_identity := app.Group("/api/v2/identity")
	{
		// Login
		api_identity.Post("/login", middlewares.MustGuest, identity.Login)
		// Logout
		api_identity.Get("/logout", middlewares.CheckRequest, identity.Logout)
		// Register
		api_identity.Post("/register", middlewares.MustGuest, identity.Register)

		// Resend email code
		api_identity.Post("/resendemail", middlewares.CheckRequest, middlewares.MustPending, identity.ReSendEmailCode)
		// Verify code
		api_identity.Post("/verifycode", middlewares.CheckRequest, middlewares.MustPending, identity.VerificationCode)
	}

	api_resource := app.Group("/api/v2/resource", middlewares.CheckRequest)
	{
		// Update Password
		api_resource.Post("/user/update/password", resource.UpdatePassword)
	}

	app.Listen(":3003")
}
