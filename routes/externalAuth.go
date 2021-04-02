package routes

import (
	"context"
	"fmt"
	"time"

	"zoom_schedule_backend_go/db"

	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"github.com/shareed2k/goth_fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExternalAuthUser struct {
	_id            primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	ExternalUserId string             `json:"externaluserid.omitempty" bson:"externaluserid,omitempty"`
	InternalUserId string             `json:"internaluserid.omitempty" bson:"internaluserid.omitempty"`
	UserName       string             `json:"name,omitempty" bson":name.omitempty"`
	Email          string             `json:"link,omitempty" bson":link.omitempty`
	Platform       string             `json:"platform.omitempty" bson:"platform.omitempty`
	AccessToken    string             `json:"accesstoken.omitempty" bson:"accesstoken.omitempty"`
	AvatarURL      string             `json:"avatarurl.omitempty" bson:"avatarurl.omitempty"`
	ExpiresAt      time.Time          `json:"expiresat.omitempty" bson:"expiresat.omitempty"`
}

const collectionExternalAuth = "externalauth"
const dbName = "zoom_schedule"

func ProviderCallback(ctx *fiber.Ctx) error {
	user, err := goth_fiber.CompleteUserAuth(ctx)
	if err != nil {
		return ctx.SendString(err.Error())
	}

	externalUser, err := GetExternalUser(ctx, user.UserID)
	if err != nil {
		return ctx.SendString(err.Error())
	}
	if externalUser == nil {
		//create User
	}

	session, _ := GetSession(ctx, user)

	return ctx.SendString(fmt.Sprintf("Welcome %v", session))
}

func GetSession(ctx *fiber.Ctx, user goth.User) (string, error) {
	store := db.GetStore()
	session, err := store.Get(ctx)
	if err != nil {
		return "", err
	}

	// Set Session Data as key/value-pair
	session.Set("userId", user.UserID)

	userId := session.Get("userId").(string)

	// save session
	if err := session.Save(); err != nil {
		return "", err
	}

	return userId, nil
}

func GetExternalUser(ctx *fiber.Ctx, externaluserID string) (*ExternalAuthUser, error) {
	collection, err := db.GetMongoDbCollection(dbName, collectionExternalAuth)
	if err != nil {
		return nil, ctx.SendString(err.Error())
	}

	ExternalAuthUser := &ExternalAuthUser{}

	err = collection.FindOne(context.TODO(), bson.D{}).Decode(&ExternalAuthUser)
	if err != nil {
		return nil, ctx.SendString(err.Error())
	}

	fmt.Println(ExternalAuthUser._id.MarshalJSON())

	return ExternalAuthUser, nil
}

func CreateUser(ctx *fiber.Ctx, user *goth.User) error {
	//collection, err := db.GetMongoDbCollection(dbName, collectionExternalAuth)
	//if err != nil {
	//	return ctx.SendString(err.Error())
	//}

	CreateInternalUser(user.NickName, user.Email)
	//if err != nil {
	//	return ctx.SendString(err.Error())
	//}

	//externalUser := &ExternalAuthUser{
	//	ExternalUserId: user.UserID,
	//	UserName:       user.NickName,
	//	Email:          user.Email,
	//	Platform:       user.Provider,
	//	AccessToken:    user.AccessToken,
	//	AvatarURL:      user.AvatarURL,
	//	ExpiresAt:      user.ExpiresAt,
	//}

	return nil
}
