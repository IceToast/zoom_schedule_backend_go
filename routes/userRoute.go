package routes

import (
	"zoom_schedule_backend_go/db"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	_id      primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	UserName string             `json:"name,omitempty" bson":name.omitempty"`
	Email    string             `json:"link,omitempty" bson":link.omitempty`
	Days     []Day              `json:"days,omitempty" bson:"days,omitempty"`
}

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

const collectionUser = "user"

func CreateInternalUser(username string, email string) (string, error) {
	collection, err := db.GetMongoDbCollection(dbName, collectionUser)
	if err != nil {
		return "", err
	}

	internalUser := &User{
		UserName: username,
		Email:    email,
	}

}
