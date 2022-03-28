package models

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// type CartProduct struct {
// 	ProductID primitive.ObjectID `json:"product_id,omitempty" bson:"product_id,omitempty"`
// 	Quantity  int64              `json:"quantity,omitempty" bson:"quantity,omitempty"`
// }

// type UserProfile struct {
// 	Fullname string `json:"fullname,omitempty" bson:"fullname,omitempty"`
// 	Age      string `json:"age,omitempty" bson:"age,omitempty"`
// 	Address  string `json:"address,omitempty" bson:"address,omitempty"`
// 	Gender   string `json:"gender,omitempty" bson:"gender,omitempty"`
// 	Phone    string `json:"phone,omitempty" bson:"phone,omitempty"`
// }

type User struct {
	mgm.DefaultModel `bson:",inline"`
	UID              string               `json:"uid,omitempty" bson:"uid,omitempty"`
	ChatIDs          []primitive.ObjectID `json:"chat_ids,omitempty" bson:"chat_ids,omitempty"`
	Email            string               `json:"email,omitempty" bson:"email,omitempty"`
	Password         string               `json:"password,omitempty" bson:"password,omitempty"`
	State            string               `json:"state,omitempty" bson:"state,omitempty"`
	Role             string               `json:"role,omitempty" bson:"role,omitempty"`
	// UserProfile      UserProfile          `json:"user_profile,omitempty" bson:"user_profile,omitempty"`
	// Cart             []CartProduct        `json:"cart,omitempty" bson:"cart,omitempty"`
}

func (model *User) Collection() *mgm.Collection {
	client, err := mgm.NewClient(options.Client().ApplyURI("mongodb://root:12345@localhost:27017"))

	if err != nil {
		panic(err)
	}

	// Get the model's db
	db := client.Database("userDB")

	// return the model's custom collection
	return mgm.NewCollection(db, "users")
}
