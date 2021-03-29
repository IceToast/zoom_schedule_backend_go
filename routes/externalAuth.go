package routes

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExternalAuth struct {
	_id            primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	ExternalUserId int                `json:"externaluserid.omitempty" bson:"externaluserid,omitempty"`
	InternalUserId string             `json:"internaluserid.omitempty" bson:"internaluserid.omitempty"`
	UserName       string             `json:"name,omitempty" bson":name.omitempty"`
	Email          string             `json:"link,omitempty" bson":link.omitempty`
	Plattform      string             `json:"plattform.omitempty" bson:"plattform.omitempty`
	AccessToken    string             `json:"accesstoken.omitempty" bson:"accesstoken.omitempty"`
	AvatarLink     string             `json:"avatarlink.omitempty" bson:"avatarlink.omitempty"`
	ExpiresAt      primitive.DateTime `json:"expiresat.omitempty" bson:"expiresat.omitempty"`
}

const dbName = "zoom_schedule"
const collectionExternalAuth = "externalauth"
