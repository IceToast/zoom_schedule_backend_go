package routes

import (
	"context"
	"encoding/json"

	"zoom_schedule_backend_go/db"

	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	_id      primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	UserName string             `json:"name,omitempty" bson":name.omitempty"`
	Email    string             `json:"link,omitempty" bson":link.omitempty`
	Days     []Day              `json:"days,omitempty" bson:"days,omitempty"`
}

type Meeting struct {
	_id      primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Name     string             `json:"name,omitempty" bson":name.omitempty"`
	Link     string             `json:"link,omitempty" bson":link.omitempty`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

type Day struct {
	_id      primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name.omitempty"`
	Meetings []Meeting          `json:"meetings,omitempty" bson:"meetings.omitempty"`
}

const collectionUser = "user"

func AuthUser(ctx *fiber.Ctx, user goth.User) error {

	return nil

}

func GetMeeting(ctx *fiber.Ctx) error {
	collection, err := db.GetMongoDbCollection(dbName, collectionUser)
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}

	var filter bson.M = bson.M{}

	if ctx.Params("id") != "" {
		id := ctx.Params("id")
		objID, _ := primitive.ObjectIDFromHex(id)
		filter = bson.M{"_id": objID}
	}

	var results []bson.M
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	defer cur.Close(context.Background())

	_ = cur.All(context.Background(), &results)

	if results == nil {
		return ctx.SendStatus(404)
	}

	jsonResults, _ := json.Marshal(results)
	return ctx.Send(jsonResults)
}
