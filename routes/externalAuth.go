package routes

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/mongodb"
	"github.com/markbates/goth"
	"github.com/shareed2k/goth_fiber"
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
const collectionSessions = "sessions"

var storage = mongodb.New(mongodb.Config{
	ConnectionURI: os.Getenv("CONNECTION_STRING"),
	Database:      dbName,
	Collection:    collectionSessions,
	Reset:         false,
})

// Storage: MongoDB, Expiration: 30days
var store = session.New(session.Config{
	Storage:    storage,
	Expiration: 720 * time.Hour,
})

func ProviderCallback(ctx *fiber.Ctx) error {
	user, err := goth_fiber.CompleteUserAuth(ctx)
	if err != nil {
		return ctx.SendString(err.Error())
	}

	name, err := GetSession(ctx, user)

	return ctx.SendString(fmt.Sprintf("Welcome %v", name))
}

func GetSession(ctx *fiber.Ctx, user goth.User) (*string, error) {
	session, err := store.Get(ctx)
	if err != nil {
		return nil, err
	}

	// Set key/value
	session.Set("name", user.Email)

	name := session.Get("name").(*string)

	// save session
	if err := session.Save(); err != nil {
		return nil, err
	}

	return name, nil
}
