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

type ProductPayload struct {
	CategoryID  string  `json:"category_id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Discount    float64 `json:"discount"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
}

func CreateProduct(c *fiber.Ctx) error {
	payload := new(ProductPayload)

	if err := c.BodyParser(payload); err != nil {
		return c.Status(422).JSON(controllers.FailedToParseBody)
	}

	id, err := primitive.ObjectIDFromHex(payload.CategoryID)

	if err != nil {
		return c.Status(500).JSON(controllers.ServerInternalError)
	}

	product := &models.Product{
		CategoryID:  id,
		Name:        payload.Name,
		Price:       payload.Price,
		Discount:    payload.Discount,
		Description: payload.Description,
		Image:       payload.Image,
	}

	if err := collection.Product.Create(product); err != nil {
		return c.Status(500).JSON(controllers.ServerInternalError)
	}

	return c.JSON(product)
}

func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(422).JSON(controllers.FailedToParseBody)
	}

	payload := new(ProductPayload)

	if err := c.BodyParser(payload); err != nil {
		return c.Status(422).JSON(controllers.FailedToParseBody)
	}

	ProductID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(500).JSON(controllers.ServerInternalError)
	}

	product := new(models.Product)
	collection.Product.FindOne(context.Background(), bson.M{"_id": ProductID}).Decode(&product)

	if product.Name == "" {
		return c.Status(500).JSON(controllers.ServerInternalError)
	}

	if payload.Description != "" {
		product.Description = payload.Description
	}

	if payload.Name != "" {
		product.Name = payload.Name
	}

	if payload.Price != 0 {
		product.Price = payload.Price
	}

	if payload.Discount != 0 {
		product.Discount = payload.Discount
	}

	if payload.Image != "" {
		product.Image = payload.Image
	}

	if payload.CategoryID != "" {
		CategoryID, err := primitive.ObjectIDFromHex(payload.CategoryID)

		if err != nil {
			return c.Status(500).JSON(controllers.ServerInternalError)
		}

		product.CategoryID = CategoryID
	}

	if err := collection.Product.Update(product); err != nil {
		return c.Status(500).JSON(controllers.ServerInternalError)
	}

	return c.JSON(product)
}

func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(422).JSON(controllers.FailedToParseBody)
	}

	ProductID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(500).JSON(controllers.ServerInternalError)
	}

	product := new(models.Product)

	collection.Product.FindOne(context.Background(), bson.M{"_id": ProductID}).Decode(&product)

	if product.Name == "" {
		return c.Status(500).JSON(controllers.ServerInternalError)
	}

	if err := collection.Product.Delete(product); err != nil {
		return c.Status(500).JSON(controllers.ServerInternalError)
	}

	return c.JSON(200)
}

type CategoryPayload struct {
	Name string `json:"name"`
}

func CreateCategory(c *fiber.Ctx) error {
	payload := new(CategoryPayload)

	if err := c.BodyParser(payload); err != nil {
		return c.Status(422).JSON(controllers.FailedToParseBody)
	}

	category := &models.Category{
		Name: payload.Name,
	}

	if err := collection.Category.Create(category); err != nil {
		return c.Status(500).JSON(controllers.ServerInternalError)
	}

	return c.JSON(category)
}

func DeleteCategory(c *fiber.Ctx) error {
	payload := c.Params("id")

	if payload == "" {
		return c.Status(500).JSON(controllers.ServerInternalError)
	}

	CategoryID, err := primitive.ObjectIDFromHex(payload)

	if err != nil {
		return c.Status(500).JSON(controllers.ServerInternalError)
	}

	category := new(models.Category)
	collection.Category.FindOne(context.Background(), bson.M{"_id": CategoryID}).Decode(&category)

	if category.Name == "" {
		return c.Status(500).JSON(controllers.ServerInternalError)
	}

	if err := collection.Category.Delete(category); err != nil {
		return c.Status(500).JSON(controllers.ServerInternalError)
	}

	return c.JSON(200)
}
