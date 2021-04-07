package routes

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExternalAuthUser struct {
	Id             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ExternalUserId string             `json:"externaluserid,omitempty" bson:"externaluserid,omitempty"`
	InternalUserId string             `json:"internaluserid,omitempty" bson:"internaluserid,omitempty"`
	UserName       string             `json:"name,omitempty" bson":name,omitempty"`
	Email          string             `json:"link,omitempty" bson":link,omitempty`
	Platform       string             `json:"platform,omitempty" bson:"platform,omitempty`
	AccessToken    string             `json:"accesstoken,omitempty" bson:"accesstoken,omitempty"`
	AvatarURL      string             `json:"avatarurl,omitempty" bson:"avatarurl,omitempty"`
	ExpiresAt      time.Time          `json:"expiresat,omitempty" bson:"expiresat,omitempty"`
}

type HTTPError struct {
	status  string
	message string
}

type User struct {
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserName string             `json:"name,omitempty" bson":name,omitempty"`
	Email    string             `json:"link,omitempty" bson":link,omitempty`
	Days     []Day              `json:"days,omitempty" bson:"days,omitempty"`
}

type Day struct {
	Name     string     `json:"name,omitempty" bson:"name,omitempty"`
	Meetings *[]Meeting `json:"meetings,omitempty" bson:"meetings,omitempty"`
}

type Meeting struct {
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson":name,omitempty"`
	Link     string             `json:"link,omitempty" bson":link,omitempty`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

type updateMeetingData struct {
	Id       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Link     string `json:"link"`
	Password string `json:"password"`
	Day      string `json:"day"`
}

type createMeetingData struct {
	Name     string `json:"name"`
	Link     string `json:"link"`
	Password string `json:"password`
	Day      string `json:"day"`
}

type deleteMeetingData struct {
	Id  string `json:"id"`
	Day string `json:"day"`
}
