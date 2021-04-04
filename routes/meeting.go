package routes

import (
	"context"
	"encoding/json"
	"zoom_schedule_backend_go/db"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Day struct {
	Name     string    `json:"name,omitempty" bson:"name,omitempty"`
	Meetings []Meeting `json:"meetings,omitempty" bson:"meetings,omitempty"`
}

type Meeting struct {
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson":name,omitempty"`
	Link     string             `json:"link,omitempty" bson":link,omitempty`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

func GetMeetings(ctx *fiber.Ctx) error {
	// Check for valid Cookie first
	store := db.GetStore()
	session, err := store.Get(ctx)
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}

	internalUserId, ok := session.Get("internalUserId").(string)

	if !ok || internalUserId == "" {
		session.Destroy()
		return ctx.Status(403).SendString("Invalid Cookie or session expired")
	}

	collection, err := db.GetMongoDbCollection(dbName, collectionUser)
	if err != nil {
		return ctx.Status(500).SendString(err.Error())

	}

	objID, _ := primitive.ObjectIDFromHex(internalUserId)

	var result *User
	err = collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return err
		}
	}

	if result.Days == nil {
		return ctx.SendString("This user has no meetings")
	}

	json, _ := json.Marshal(result.Days)
	ctx.Send(json)
	return nil
}

func CreateMeeting(ctx *fiber.Ctx) error {

	collection, err := db.GetMongoDbCollection(dbName, collectionUser)
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
	collection, err := db.GetMongoDbCollection(dbName, collectionUser)
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
	collection, err := db.GetMongoDbCollection(dbName, collectionUser)
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
