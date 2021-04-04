package routes

import (
	"context"
	"encoding/json"
	"zoom_schedule_backend_go/db"
	"zoom_schedule_backend_go/helpers"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HTTPError struct {
	status  string
	message string
}

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

// GetMeetings godoc
// @Summary Retrieves meetings from the local Mongo database for a certain user.
// @Description Resolves a userId via a given session cookie. The backend throws an error if the cookie does not exist.
// @Accept json
// @Produce json
// @Success 200 {object} Meeting
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /api/meeting [get]
func GetMeetings(ctx *fiber.Ctx) error {
	internalUserId, err := helpers.VerifyCookie(ctx)
	if err != nil {
		return ctx.Status(403).SendString(err.Error())
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


// CreateMeeting godoc
// @Summary Creates a meeting in the local Mongo database.
// @Description Requires a JSON encoded Meeting object in the body.
// @Accept json
// @Produce json
// @Success 200
// @Failure 500 {object} HTTPError
// @Router /api/meeting [post]
func CreateMeeting(ctx *fiber.Ctx) error {
	//Verify Cookie
	//internalUserId, err := helpers.VerifyCookie(ctx)
	//if err != nil {
	//	return ctx.Status(403).SendString(err.Error())
	//}

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

// UpdateMeeting godoc
// @Summary Updates a meeting in the local Mongo database.
// @Description Requires a userId
// @Accept json
// @Produce json
// @Success 200
// @Failure 500 {object} HTTPError
// @Router /api/meeting [put]
func UpdateMeeting(ctx *fiber.Ctx) error {
	//Verify Cookie
	internalUserId, err := helpers.VerifyCookie(ctx)
	if err != nil {
		return ctx.Status(403).SendString(err.Error())
	}

	collection, err := db.GetMongoDbCollection(dbName, collectionUser)
	if err != nil {
		ctx.Status(500).SendString(err.Error())
		return err
	}

	var meeting Meeting
	json.Unmarshal(ctx.Body(), &meeting)

	update := bson.M{
		"$set": meeting,
	}

	objID, _ := primitive.ObjectIDFromHex(internalUserId)
	res, err := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)

	if err != nil {
		ctx.Status(500).SendString(err.Error())
		return err
	}

	response, _ := json.Marshal(res)
	ctx.Send(response)
	return nil
}


// DeleteMeeting godoc
// @Summary Deletes a meeting in the local Mongo database.
// @Description Requires a userId
// @Accept json
// @Produce json
// @Success 200
// @Failure 500 {object} HTTPError
// @Router /api/meeting [delete]
func DeleteMeeting(ctx *fiber.Ctx) error {
	//Verify Cookie
	internalUserId, err := helpers.VerifyCookie(ctx)
	if err != nil {
		return ctx.Status(403).SendString(err.Error())
	}

	collection, err := db.GetMongoDbCollection(dbName, collectionUser)

	objID, _ := primitive.ObjectIDFromHex(internalUserId)
	res, err := collection.DeleteOne(context.Background(), bson.M{"_id": objID})

	if err != nil {
		ctx.Status(500).SendString(err.Error())
		return err
	}

	jsonResponse, _ := json.Marshal(res)
	ctx.Send(jsonResponse)
	return nil
}
