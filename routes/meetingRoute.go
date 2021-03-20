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

func GetMeeting(c *fiber.Ctx) error {
	collection, err := db.GetMongoDbCollection(dbName, collectionName)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	var filter bson.M = bson.M{}

	if c.Params("id") != "" {
		id := c.Params("id")
		objID, _ := primitive.ObjectIDFromHex(id)
		filter = bson.M{"_id": objID}
	}

	var results []bson.M
	cur, err := collection.Find(context.Background(), filter)
	defer cur.Close(context.Background())

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	_ = cur.All(context.Background(), &results)

	if results == nil {
		return c.SendStatus(404)
	}

	jsonResults, _ := json.Marshal(results)
	return c.Send(jsonResults)
}

func CreateMeeting(c *fiber.Ctx) error {

	collection, err := db.GetMongoDbCollection(dbName, collectionName)
	if err != nil {

		return c.Status(500).SendString(err.Error())
	}

	var meeting Meeting
	_ = json.Unmarshal(c.Body(), &meeting)

	res, err := collection.InsertOne(context.Background(), meeting)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	response, _ := json.Marshal(res)
	return c.Send(response)
}

func UpdateMeeting(c *fiber.Ctx) error {
	collection, err := db.GetMongoDbCollection(dbName, collectionName)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	var meeting Meeting
	_ = json.Unmarshal(c.Body(), &meeting)

	update := bson.M{
		"$set": meeting,
	}

	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	res, err := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	response, _ := json.Marshal(res)
	return c.Send(response)
}

func DeleteMeeting(c *fiber.Ctx) error {
	collection, err := db.GetMongoDbCollection(dbName, collectionName)

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	res, err := collection.DeleteOne(context.Background(), bson.M{"_id": objID})

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	jsonResponse, _ := json.Marshal(res)
	return c.Send(jsonResponse)
}
