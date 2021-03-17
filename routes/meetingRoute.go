package routes

import (
	"context"
	"encoding/json"

	"zoom_schedule_backend_go/db"

	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Meeting struct {
	_id primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Name	string `json:"name,omitempty" bson":name.omitempty"`
	Link	string `json:"link,omitempty" bson":link.omitempty`
	Password	string `json:"password,omitempty bson:"password,omitempty"`
}

const dbName = "zoom_schedule"
const collectionName = "meeting"

func GetMeeting(c *fiber.Ctx) {
	collection, err := db.GetMongoDbCollection(dbName, collectionName)
	if err != nil {
		c.Status(500).Send(err)
		return
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
		c.Status(500).Send(err)
		return
	}

	cur.All(context.Background(), &results)

	if results == nil {
		c.SendStatus(404)
		return
	}

	json, _ := json.Marshal(results)
	c.Send(json)
}

func CreateMeeting(c *fiber.Ctx) {

	collection, err := db.GetMongoDbCollection(dbName, collectionName)
	if err != nil {
		c.Status(500).Send(err)
		return
	}

	var meeting Meeting
	json.Unmarshal([]byte(c.Body()), &meeting)

	res, err := collection.InsertOne(context.Background(), meeting)
	if err != nil {
		c.Status(500).Send(err)
		return
	}

	response, _ := json.Marshal(res)
	c.Send(response)
}

func UpdateMeeting(c *fiber.Ctx) {
	collection, err := db.GetMongoDbCollection(dbName, collectionName)
	if err != nil {
		c.Status(500).Send(err)
		return
	}
	var meeting Meeting
	json.Unmarshal([]byte(c.Body()), &meeting)

	update := bson.M{
		"$set": meeting,
	}

	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	res, err := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)

	if err != nil {
		c.Status(500).Send(err)
		return
	}

	response, _ := json.Marshal(res)
	c.Send(response)
}

func DeleteMeeting(c *fiber.Ctx) {
	collection, err := db.GetMongoDbCollection(dbName, collectionName)

	if err != nil {
		c.Status(500).Send(err)
		return
	}

	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	res, err := collection.DeleteOne(context.Background(), bson.M{"_id": objID})

	if err != nil {
		c.Status(500).Send(err)
		return
	}

	jsonResponse, _ := json.Marshal(res)
	c.Send(jsonResponse)
}