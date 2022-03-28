package admin

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/shinhagunn/Shop-Watches/backend/config/collection"
	"github.com/shinhagunn/Shop-Watches/backend/controllers"
	"github.com/shinhagunn/Shop-Watches/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetComments(c *fiber.Ctx) error {
	comments := []models.Comment{}

	if err := collection.Comment.SimpleFind(&comments, bson.M{}); err != nil {
		return c.Status(422).JSON(controllers.FailedConnectDataInDatabase)
	}

	return c.JSON(comments)
}

func DeleteComment(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(422).JSON(controllers.FailedToParseBody)
	}

	CommentID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(500).JSON(controllers.ServerInternalError)
	}

	comment := new(models.Comment)

	collection.Comment.FindOne(context.Background(), bson.M{"_id": CommentID}).Decode(&comment)

	if comment.Content == "" {
		return c.Status(500).JSON(controllers.ServerInternalError)
	}

	if err := collection.Comment.Delete(comment); err != nil {
		return c.Status(500).JSON(controllers.ServerInternalError)
	}

	return c.JSON(200)
}
