package routes

import (
	"context"
	"encoding/json"
	"zoom_schedule_backend_go/db"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Day struct {
	_id      primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name.omitempty"`
	Meetings []Meeting          `json:"meetings,omitempty" bson:"meetings.omitempty"`
}

type Meeting struct {
	_id      primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Name     string             `json:"name,omitempty" bson":name.omitempty"`
	Link     string             `json:"link,omitempty" bson":link.omitempty`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

const collectionName = "meeting"

func GetMeeting(ctx *fiber.Ctx) error {
	collection, collectionError := db.GetMongoDbCollection(dbName, collectionName)
	store := db.GetStore()
	session, storeError := store.Get(ctx)

	if collectionError != nil {
		ctx.Status(500).SendString(collectionError.Error())
		return collectionError
	}

	if storeError != nil {
		ctx.Status(500).SendString(storeError.Error())
		return storeError
	}

	var filter = bson.M{}
	userId := session.Get("userId").(string)

	if userId != "" {
		objID, _ := primitive.ObjectIDFromHex(userId)
		filter = bson.M{"_id": objID}
	}

	var results []bson.M
	cur, err := collection.Find(context.Background(), filter)
	defer cur.Close(context.Background())

	if err != nil {
		ctx.Status(500).SendString(err.Error())
		return err
	}

	cur.All(context.Background(), &results)

	if results == nil {
		ctx.SendStatus(404)
		return nil
	}

	json, _ := json.Marshal(results)
	ctx.Send(json)
	return nil
}

func CreateMeeting(ctx *fiber.Ctx) error {

	collection, err := db.GetMongoDbCollection(dbName, collectionName)
	if err != nil {
		ctx.Status(500).SendString(err.Error())
		return err
	}

	var meeting Meeting
	json.Unmarshal(ctx.Body(), &meeting)

	res, err := collection.InsertOne(context.Background(), meeting)
	if err != nil {
		ctx.Status(500).SendString(err.Error())
		return err
	}

	response, _ := json.Marshal(res)
	ctx.Send(response)
	return nil
}

func UpdateMeeting(ctx *fiber.Ctx) error {
	collection, err := db.GetMongoDbCollection(dbName, collectionName)
	store := db.GetStore()
	session, storeError := store.Get(ctx)

	if err != nil {
		ctx.Status(500).SendString(err.Error())
		return err
	}

	if storeError != nil {
		ctx.Status(500).SendString(storeError.Error())
		return err
	}

	var meeting Meeting
	json.Unmarshal(ctx.Body(), &meeting)

	update := bson.M{
		"$set": meeting,
	}

	objID, _ := primitive.ObjectIDFromHex(session.Get("userId").(string))
	res, err := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)

	if err != nil {
		ctx.Status(500).SendString(err.Error())
		return err
	}

	response, _ := json.Marshal(res)
	ctx.Send(response)
	return nil
}

func DeleteMeeting(ctx *fiber.Ctx) error {
	collection, err := db.GetMongoDbCollection(dbName, collectionName)
	store := db.GetStore()
	session, storeError := store.Get(ctx)

	if err != nil {
		ctx.Status(500).SendString(err.Error())
		return err
	}

	if storeError != nil {
		ctx.Status(500).SendString(storeError.Error())
		return err
	}

	objID, _ := primitive.ObjectIDFromHex(session.Get("userId").(string))
	res, err := collection.DeleteOne(context.Background(), bson.M{"_id": objID})

	if err != nil {
		ctx.Status(500).SendString(err.Error())
		return err
	}

	jsonResponse, _ := json.Marshal(res)
	ctx.Send(jsonResponse)
	return nil
}
