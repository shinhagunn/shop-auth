package public

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shinhagunn/Shop-Watches/backend/config/collection"
	"github.com/shinhagunn/Shop-Watches/backend/controllers"
	"github.com/shinhagunn/Shop-Watches/backend/models"
	"go.mongodb.org/mongo-driver/bson"
)

func GetSlides(c *fiber.Ctx) error {
	slides := []models.Slide{}

	if err := collection.Slide.SimpleFind(&slides, bson.M{}); err != nil {
		return c.Status(422).JSON(controllers.FailedConnectDataInDatabase)
	}

	return c.JSON(slides)
}
