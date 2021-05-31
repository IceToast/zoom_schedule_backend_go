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
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetMeetings godoc
// @Summary Retrieves all meetings from the database for a certain user.
// @Description Resolves a userId via a given session cookie. The backend throws an error if the cookie does not exist.
// @Accept json
// @Produce json
// @Success 200 {object} []Day
// @Failure 403 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /api/meeting [get]
func GetMeetings(ctx *fiber.Ctx) error {
	//Verify Cookie
	internalUserId, err := helpers.VerifyCookie(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).SendString(err.Error())
	}

	collection, err := db.GetMongoDbCollection(dbName, collectionUser)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	objID, _ := primitive.ObjectIDFromHex(internalUserId)

	var result *User
	err = collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return ctx.Status(fiber.StatusInternalServerError).SendString("User not found")
		}
	}

	db.CloseMongoDbConnection(collection)

	if result.Days == nil {
		return ctx.SendString("This user has no meetings")
	}

	json, _ := json.Marshal(result.Days)
	return ctx.Status(fiber.StatusOK).Send(json)

}

// CreateMeeting godoc
// @Summary Creates a meeting in the database.
// @Description Requires a JSON encoded Meeting object in the body.
// @Param request body createMeetingData true "Meeting Data required to create a Meeting"
// @Accept json
// @Produce json
// @Success 200 {object} Meeting
// @Failure 403	{object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /api/meeting [post]
func CreateMeeting(ctx *fiber.Ctx) error {
	//Verify Cookie
	internalUserId, err := helpers.VerifyCookie(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	collection, err := db.GetMongoDbCollection(dbName, collectionUser)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())

	}

	var meetingData createMeetingData
	//Convert HTTP POST Data to Struct
	json.Unmarshal(ctx.Body(), &meetingData)

	//Do not update if request Data is invalid -> no Day to select or no meeting properties
	if meetingData.Name == "" && meetingData.Link == "" && meetingData.Password == "" || meetingData.Day == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("No valid Meeting")
	}

	//Convert Ids of type string to type "ObjectIds"
	userObjID, _ := primitive.ObjectIDFromHex(internalUserId)
	meetingObjId := primitive.NewObjectID()

	//Map request data to meeting update struct
	update := bson.M{
		"$push": bson.M{
			"days.$.meetings": Meeting{
				Id:        meetingObjId,
				Name:      meetingData.Name,
				Link:      meetingData.Link,
				Password:  meetingData.Password,
				StartTime: meetingData.StartTime,
				EndTime:   meetingData.EndTime,
			},
		},
	}

	//Filter User and Day by Meeting Data from request
	filter := bson.M{"_id": userObjID, "days.name": meetingData.Day}

	res, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if res.ModifiedCount < 1 {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Could not create Meeting")
	}

	db.CloseMongoDbConnection(collection)

	meeting := Meeting{
		Id:        meetingObjId,
		Name:      meetingData.Name,
		Link:      meetingData.Link,
		Password:  meetingData.Password,
		StartTime: meetingData.StartTime,
	}

	response, _ := json.Marshal(meeting)
	return ctx.Status(fiber.StatusOK).Send(response)
}

// UpdateMeeting godoc
// @Summary Updates a meeting in the database.
// @Description Requires a JSON encoded Meeting object in the body
// @Accept json
// @Produce json
// @Param request body updateMeetingData true "Meeting Data required to update a Meeting"
// @Success 200
// @Failure 403	{object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /api/meeting [put]
func UpdateMeeting(ctx *fiber.Ctx) error {
	//Verify Cookie
	internalUserId, err := helpers.VerifyCookie(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	collection, err := db.GetMongoDbCollection(dbName, collectionUser)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())

	}

	var meetingData updateMeetingData
	//Convert HTTP POST Data to Struct
	json.Unmarshal(ctx.Body(), &meetingData)

	//Do not update if request Data is invalid -> no Day to select or no meeting properties
	if meetingData.Name == "" && meetingData.Link == "" && meetingData.Password == "" || meetingData.Day == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("No valid Meeting")
	}

	//Convert Ids of type string to type "ObjectIds"
	userObjId, _ := primitive.ObjectIDFromHex(internalUserId)
	meetingObjId, _ := primitive.ObjectIDFromHex(meetingData.Id)

	//Map request data to meeting update struct
	update := bson.M{
		"$set": bson.M{
			"days.$.meetings.$[elem]": Meeting{
				Id:        meetingObjId,
				Name:      meetingData.Name,
				Link:      meetingData.Link,
				Password:  meetingData.Password,
				StartTime: meetingData.StartTime,
				EndTime:   meetingData.EndTime,
			},
		},
	}

	//Filter days array to update meeting in correct day
	filter := bson.D{{"_id", userObjId}, {"days.name", meetingData.Day}}
	//Filter meetings Array to update correct meeting
	filterArray := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{bson.M{"elem._id": meetingObjId}},
	})

	res, err := collection.UpdateOne(context.Background(), filter, update, filterArray)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if res.ModifiedCount < 1 {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Could not update Meeting")
	}
	db.CloseMongoDbConnection(collection)

	meeting := Meeting{
		Id:        meetingObjId,
		Name:      meetingData.Name,
		Link:      meetingData.Link,
		Password:  meetingData.Password,
		StartTime: meetingData.StartTime,
		EndTime:   meetingData.EndTime,
	}

	response, _ := json.Marshal(meeting)

	return ctx.Status(fiber.StatusOK).Send(response)
}

// DeleteMeeting godoc
// @Summary Deletes a meeting in the Database
// @Description Requires a JSON encoded Meeting object in the body
// @Accept json
// @Produce json
// @Param request body deleteMeetingData true "Meeting Data required to delete a Meeting"
// @Success 200
// @Failure 403	{object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /api/meeting [delete]
func DeleteMeeting(ctx *fiber.Ctx) error {
	//Verify Cookie
	internalUserId, err := helpers.VerifyCookie(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	collection, err := db.GetMongoDbCollection(dbName, collectionUser)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	var meetingData deleteMeetingData
	//Convert HTTP POST Data to Struct
	json.Unmarshal(ctx.Body(), &meetingData)

	//Do not update if request Data is invalid -> no Day to select or no meeting properties
	if meetingData.Id == "" || meetingData.Day == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("No valid Meeting")
	}

	//Convert Ids of type string to type "ObjectIds"
	userObjId, _ := primitive.ObjectIDFromHex(internalUserId)
	meetingObjId, _ := primitive.ObjectIDFromHex(meetingData.Id)

	delete := bson.M{
		"$pull": bson.M{
			"days.$.meetings": bson.M{
				"_id": meetingObjId,
			},
		},
	}

	//Filter days array to update meeting in correct day
	filter := bson.D{{"_id", userObjId}, {"days.name", meetingData.Day}}

	//Filter meetings Array to update correct meeting
	res, err := collection.UpdateOne(context.Background(), filter, delete)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if res.ModifiedCount < 1 {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Could not delete Meeting")
	}
	db.CloseMongoDbConnection(collection)

	return ctx.SendStatus(fiber.StatusOK)
}

// FlushSchedule godoc
// @Summary Deletes all Meetings of a User
// @Description Clear the Users Zoom Schedule
// @Success 200
// @Failure 403	{object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /api/meeting/flushSchedule [delete]
func FlushSchedule(ctx *fiber.Ctx) error {
	//Verify Cookie
	internalUserId, err := helpers.VerifyCookie(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	collection, err := db.GetMongoDbCollection(dbName, collectionUser)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	userObjId, _ := primitive.ObjectIDFromHex(internalUserId)
	filter := bson.D{{"_id", userObjId}}

	//Map request data to meeting update struct
	emptySchedule := bson.M{
		"$set": bson.M{
			"days": []Day{
				0: {Name: "Monday", Meetings: &[]Meeting{}},
				1: {Name: "Tuesday", Meetings: &[]Meeting{}},
				2: {Name: "Wednesday", Meetings: &[]Meeting{}},
				3: {Name: "Thursday", Meetings: &[]Meeting{}},
				4: {Name: "Friday", Meetings: &[]Meeting{}},
				5: {Name: "Saturday", Meetings: &[]Meeting{}}},
		},
	}

	res, err := collection.UpdateOne(context.Background(), filter, emptySchedule)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if res.ModifiedCount < 1 {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Could not flush Schedule")
	}
	db.CloseMongoDbConnection(collection)

	return ctx.SendStatus(fiber.StatusOK)
}
