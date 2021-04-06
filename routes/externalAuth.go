package routes

import (
	"encoding/json"
	"fmt"
	"time"

	"zoom_schedule_backend_go/db"

	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"github.com/shareed2k/goth_fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
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

const collectionExternalAuth = "externalauth"
const dbName = "zoom_schedule"

// ProviderCallback godoc
// @Summary Handles the OAuth2 authentication callback for a certain goth provider.
// @Description Parses the Fiber context to receive the user's ID and creates the user if it does not exist yet.
// @Accept json
// @Produce json
// @Success 200
// @Param provider path string true "goth_fiber provider (google/discord/github)"
// @Failure 500 {object} HTTPError
// @Router /api/{provider}/callback [get]
func ProviderCallback(ctx *fiber.Ctx) error {
	user, err := goth_fiber.CompleteUserAuth(ctx)
	if err != nil {
		return ctx.SendString(err.Error())
	}

	externalUser, err := GetExternalUser(user.UserID)
	if err == mongo.ErrNoDocuments {
		externalUser, err = CreateUser(ctx, user)
		if err != nil {
			return ctx.SendString("User creation Failed")
		}
	}

	session, _ := GetSession(ctx, externalUser)
	marshalled, _ := json.Marshal(externalUser)

	return ctx.SendString(fmt.Sprintf("Welcome %v", string(marshalled)+"Session: "+session))
}

func GetSession(ctx *fiber.Ctx, externalUser *ExternalAuthUser) (string, error) {
	store := db.GetStore()
	session, err := store.Get(ctx)
	if err != nil {
		return "", err
	}

	// Set Session Data as key/value-pair
	session.Set("internalUserId", externalUser.InternalUserId)

	userId := session.Get("internalUserId").(string)

	// save session
	if err := session.Save(); err != nil {
		return "", err
	}

	return userId, nil
}

func GetExternalUser(externaluserID string) (*ExternalAuthUser, error) {
	collection, err := db.GetMongoDbCollection(dbName, collectionExternalAuth)
	if err != nil {
		return nil, err
	}

	var result *ExternalAuthUser
	err = collection.FindOne(context.TODO(), bson.M{"externaluserid": externaluserID}).Decode(&result)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
	}

	return result, nil
}

func CreateUser(ctx *fiber.Ctx, user goth.User) (*ExternalAuthUser, error) {
	collection, err := db.GetMongoDbCollection(dbName, collectionExternalAuth)
	if err != nil {
		return nil, err
	}

	internalUserId, err := CreateInternalUser(user.Name, user.Email)
	if err != nil {
		return nil, err
	}

	externalUser := &ExternalAuthUser{
		ExternalUserId: user.UserID,
		InternalUserId: internalUserId,
		UserName:       user.Name,
		Email:          user.Email,
		Platform:       user.Provider,
		AccessToken:    user.AccessToken,
		AvatarURL:      user.AvatarURL,
		ExpiresAt:      user.ExpiresAt,
	}

	res, err := collection.InsertOne(context.Background(), externalUser)
	if err != nil {
		return nil, err
	}

	response, _ := json.Marshal(res)

	ctx.SendString(string(response))

	return externalUser, nil
}
