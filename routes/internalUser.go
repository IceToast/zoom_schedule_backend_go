package routes

import (
	"context"
	"zoom_schedule_backend_go/db"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserName string             `json:"name,omitempty" bson":name,omitempty"`
	Email    string             `json:"link,omitempty" bson":link,omitempty`
	Days     []Day              `json:"days,omitempty" bson:"days,omitempty"`
}

const collectionUser = "user"

func CreateInternalUser(username string, email string) (string, error) {
	collection, err := db.GetMongoDbCollection(dbName, collectionUser)
	if err != nil {
		return "", err
	}

	internalUser := &User{
		UserName: username,
		Email:    email,
		Days: []Day{
			0: {Name: "Monday"},
			1: {Name: "Tuesday"},
			2: {Name: "Wedensday"},
			3: {Name: "Thursday"},
			4: {Name: "Friday"},
			5: {Name: "Saturday"},
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
