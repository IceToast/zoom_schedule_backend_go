package routes

import (
	"context"
	"encoding/json"

	"zoom_schedule_backend_go/db"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Meeting struct {
	_id      primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Name     string             `json:"name,omitempty" bson":name.omitempty"`
	Link     string             `json:"link,omitempty" bson":link.omitempty`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

const dbName = "zoom_schedule"
const collectionName = "meeting"

func GetMeeting(ctx *fiber.Ctx) error {
	collection, err := db.GetMongoDbCollection(dbName, collectionName)
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
	defer cur.Close(context.Background())

	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}

	_ = cur.All(context.Background(), &results)

	if results == nil {
		return ctx.SendStatus(404)
	}

	jsonResults, _ := json.Marshal(results)
	return ctx.Send(jsonResults)
}

func CreateMeeting(ctx *fiber.Ctx) error {

	collection, err := db.GetMongoDbCollection(dbName, collectionName)
	if err != nil {

		return ctx.Status(500).SendString(err.Error())
	}

	var meeting Meeting
	_ = json.Unmarshal(ctx.Body(), &meeting)

	res, err := collection.InsertOne(context.Background(), meeting)
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}

	response, _ := json.Marshal(res)
	return ctx.Send(response)
}

func UpdateMeeting(ctx *fiber.Ctx) error {
	collection, err := db.GetMongoDbCollection(dbName, collectionName)
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	var meeting Meeting
	_ = json.Unmarshal(ctx.Body(), &meeting)

	update := bson.M{
		"$set": meeting,
	}

	objID, _ := primitive.ObjectIDFromHex(ctx.Params("id"))
	res, err := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)

	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}

	response, _ := json.Marshal(res)
	return ctx.Send(response)
}

func DeleteMeeting(ctx *fiber.Ctx) error {
	collection, err := db.GetMongoDbCollection(dbName, collectionName)

	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}

	objID, _ := primitive.ObjectIDFromHex(ctx.Params("id"))
	res, err := collection.DeleteOne(context.Background(), bson.M{"_id": objID})

	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}

	jsonResponse, _ := json.Marshal(res)
	return ctx.Send(jsonResponse)
}
