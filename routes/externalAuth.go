package routes

import (
	"fmt"
	"os"
	"time"

	"zoom_schedule_backend_go/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/mongodb"
	"github.com/markbates/goth"
	"github.com/shareed2k/goth_fiber"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExternalAuthUser struct {
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

const collectionExternalAuth = "externalauth"


func ProviderCallback(ctx *fiber.Ctx) error {
	user, err := goth_fiber.CompleteUserAuth(ctx)
	if err != nil {
		return ctx.SendString(err.Error())
	}

	session, _ := GetSession(ctx, user)

	return ctx.SendString(fmt.Sprintf("Welcome %v", session))
}

func GetSession(ctx *fiber.Ctx, user ExternalAuthUser) (string, error) {
	store := db.GetStore()
	session, err := store.Get(ctx)
	if err != nil {
		return "", err
	}

	// Set Session Data as key/value-pair
	session.Set("userId", user.InternalUserId)

	userId := session.Get("userId").(string)

	// save session
	if err := session.Save(); err != nil {
		return "", err
	}

	return name + userId, nil
}
