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
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson":name,omitempty"`
	Link      string             `json:"link,omitempty" bson":link,omitempty`
	Password  string             `json:"password,omitempty" bson:"password,omitempty"`
	StartTime string             `json:"startTime,omitempty" bson:"startTime,omitempty"`
	EndTime   string             `json:"endTime,omitempty" bson:"endTime,omitempty"`
}

type updateMeetingData struct {
	Id        string `json:"id,omitempty validate:"required"`
	Name      string `json:"name"`
	Link      string `json:"link"`
	Password  string `json:"password"`
	Day       string `json:"day" validate:"required"`
	StartTime string `json:"startTime,omitempty"`
	EndTime   string `json:"endTime,omitempty"`
}

type createMeetingData struct {
	Name      string `json:"name" validate:"required"`
	Link      string `json:"link"`
	Password  string `json:"password`
	Day       string `json:"day" validate:"required"`
	StartTime string `json:"startTime,omitempty"`
	EndTime   string `json:"endTime,omitempty"`
}

type deleteMeetingData struct {
	Id  string `json:"id" validate:"required"`
	Day string `json:"day" validate:"required"`
}

type userData struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	AvatarUrl string `json:"avatarurl"`
	Platform  string `json:"platform"`
}
