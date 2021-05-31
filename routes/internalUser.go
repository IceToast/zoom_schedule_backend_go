package routes

import (
	"context"
	"zoom_schedule_backend_go/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const collectionUser = "user"

func CreateInternalUser(username string, email string) (string, error) {
	collection := db.DbInstance.DB.Collection(collectionUser)

	internalUser := &User{
		UserName: username,
		Email:    email,
		Days: []Day{
			0: {Name: "Monday", Meetings: &[]Meeting{}},
			1: {Name: "Tuesday", Meetings: &[]Meeting{}},
			2: {Name: "Wednesday", Meetings: &[]Meeting{}},
			3: {Name: "Thursday", Meetings: &[]Meeting{}},
			4: {Name: "Friday", Meetings: &[]Meeting{}},
			5: {Name: "Saturday", Meetings: &[]Meeting{}},
		},
	}

	res, err := collection.InsertOne(context.Background(), internalUser)
	if err != nil {
		return "", err
	}

	internalUserId, _ := res.InsertedID.(primitive.ObjectID)
	internalUserIdString := internalUserId.Hex()

	return internalUserIdString, nil
}

func DeleteInternalUser(internalUserId string) error {
	collection := db.DbInstance.DB.Collection(collectionUser)

	userObjId, _ := primitive.ObjectIDFromHex(internalUserId)

	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": userObjId})
	if err != nil {
		return err
	}

	return nil
}
