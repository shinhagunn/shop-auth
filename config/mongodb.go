package config

import (
	"context"
	"fmt"

	"github.com/kamva/mgm/v3"
	"github.com/shinhagunn/Shop-Watches/backend/config/collection"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongoDB() {
	mgm.SetDefaultConfig(&mgm.Config{}, "mongodb")
	client, err := mgm.NewClient(options.Client().ApplyURI("mongodb://root:123456@localhost:27017"))

	if err != nil {
		panic(err)
	}

	InitProductDB(client)
	InitUserDB(client)
	InitOrderDB(client)
}

func InitProductDB(client *mongo.Client) {
	// client, err := mgm.SetDefaultConfig()
	// mgm.Coll()

	collection.Product = mgm.NewCollection(client.Database("productDB"), "products")
	collection.Category = mgm.NewCollection(client.Database("productDB"), "categories")
	collection.Slide = mgm.NewCollection(client.Database("productDB"), "slides")

	fmt.Println("Connected to ProductDB!")

	collection.Category.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.M{"name": 1},
		Options: options.Index().SetUnique(true),
	})
}

func InitUserDB(client *mongo.Client) {
	collection.User = mgm.NewCollection(client.Database("userDB"), "users")
	collection.Code = mgm.NewCollection(client.Database("userDB"), "codes")
	collection.Custommer = mgm.NewCollection(client.Database("userDB"), "custommers")

	fmt.Println("Connected to UserDB!")

	collection.User.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"uid": 1},
			Options: options.Index().SetUnique(true),
		},
	})
}

func InitOrderDB(client *mongo.Client) {
	collection.Order = mgm.NewCollection(client.Database("orderDB"), "orders")

	// collection.Order = mgm.Coll(&models.Order{})

	fmt.Println("Connected to OrderDB!")
}
