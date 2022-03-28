package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/shinhagunn/Shop-Watches/backend/config"
	"github.com/shinhagunn/Shop-Watches/backend/config/collection"
	"github.com/shinhagunn/Shop-Watches/backend/controllers"
	"github.com/shinhagunn/Shop-Watches/backend/controllers/admin"
	"github.com/shinhagunn/Shop-Watches/backend/controllers/identity"
	"github.com/shinhagunn/Shop-Watches/backend/controllers/public"
	"github.com/shinhagunn/Shop-Watches/backend/controllers/resource"
	"github.com/shinhagunn/Shop-Watches/backend/models"
	"github.com/shinhagunn/Shop-Watches/backend/routes/middlewares"
	"github.com/shinhagunn/Shop-Watches/backend/utils"
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
		api_identity.Get("/logout", middlewares.CheckJWT, identity.Logout)
		// Register
		api_identity.Post("/register", middlewares.MustGuest, identity.Register)

		// Resend email code
		api_identity.Post("/resendemail", middlewares.CheckJWT, middlewares.MustPending, identity.ReSendEmailCode)
		// Verify code
		api_identity.Post("/verifycode", middlewares.CheckJWT, middlewares.MustPending, identity.VerificationCode)
	}

	api_public := app.Group("/api/v2/public")
	{
		// Get all slides
		api_public.Get("/slides", public.GetSlides)
		// Get products
		api_public.Get("/products", public.GetProducts)
		// Get product by id
		api_public.Get("/product/:id", public.GetProductByID)
		// Get categories
		api_public.Get("/categories", public.GetCategories)
		// Post message custommer
		api_public.Post("/sendmessage", public.CustommerMessage)
	}

	api_resource := app.Group("/api/v2/resource", middlewares.CheckJWT)
	{
		// Get user info
		api_resource.Get("/user", middlewares.MustActive, resource.GetUserCurrent)
		// Update password
		api_resource.Post("/user/update/password", middlewares.MustActive, resource.UpdatePassword)
		// Update profile
		api_resource.Post("/user/update/profile", middlewares.MustActive, resource.UpdateUserProfile)
		// Add product to cart
		api_resource.Post("/user/cart", middlewares.MustActive, resource.AddProductToCart)
		// Get all product in cart
		api_resource.Get("/user/cart", middlewares.MustActive, resource.GetCartProducts)
		// Remove product in cart
		api_resource.Delete("/user/cart/:id", middlewares.MustActive, resource.RemoveProductInCart)
		// Get comments by product
		api_public.Get("/product/:id/comments", middlewares.MustActive, resource.GetCommentsInProduct)
		// Add comment
		api_resource.Post("/product/:id/comment", middlewares.MustActive, resource.CreateComment)
		// Like comment
		api_resource.Get("/comment/:id/like", middlewares.MustActive, resource.LikeComment)
		// Dislike comment
		api_resource.Get("/comment/:id/dislike", middlewares.MustActive, resource.DislikeComment)
		// Add order
		api_resource.Get("/user/order", middlewares.MustActive, resource.HandleOrder)
	}

	api_admin := app.Group("/api/v2/admin", middlewares.CheckJWT)
	{
		// Add category
		api_admin.Post("/category", middlewares.MustAdmin, admin.CreateCategory)
		// Delete category
		api_admin.Delete("/category/:id", middlewares.MustAdmin, admin.DeleteCategory)
		// Add product
		api_admin.Post("/product", middlewares.MustAdmin, admin.CreateProduct)
		// Update product
		api_admin.Post("/product/:id", middlewares.MustAdmin, admin.UpdateProduct)
		// Delete product
		api_admin.Delete("/product/:id", middlewares.MustAdmin, admin.DeleteProduct)
		// Get comments
		api_admin.Get("/comments", middlewares.MustAdmin, admin.GetComments)
		// Remove comment
		api_admin.Delete("/comment/:id", middlewares.MustAdmin, admin.DeleteComment)
		// Add slide
		api_admin.Post("/slide", middlewares.MustAdmin, admin.CreateSlide)
		// Update slide
		api_admin.Post("/slide/:id", middlewares.MustAdmin, admin.UpdateSlide)
		// Delete slide
		api_admin.Delete("/slide/:id", middlewares.MustAdmin, admin.DeleteSlide)
	}

	app.Listen(":3003")
}
