package routes

import (
	"encoding/json"

	"zoom_schedule_backend_go/db"

	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"github.com/shareed2k/goth_fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

const (
	collectionExternalAuth = "externalauth"
	dbName                 = "zoom_schedule"
	webAppUrl              = "https://zoom.icetoast.cloud"
)

// ProviderCallback godoc
// @Summary Handles the OAuth2 authentication callback for a certain goth provider.
// @Description Parses the Fiber context to receive the user's ID and creates the user if it does not exist yet.
// @Accept json
// @Produce json
// @Success 200
// @Param provider path string true "Google, Discord"
// @Failure 500 {object} HTTPError
// @Router /api/{provider}/callback [get]
func ProviderCallback(ctx *fiber.Ctx) error {
	user, err := goth_fiber.CompleteUserAuth(ctx)
	if err != nil {
		return ctx.SendString(err.Error())
	}

	var externalUser *ExternalAuthUser
	externalUser, err = GetExternalUser(user.UserID)
	if err == mongo.ErrNoDocuments {
		externalUser, err = CreateUser(ctx, user)
		if err != nil {
			return ctx.SendString("User creation Failed")
		}
	}

	GetSession(ctx, externalUser)

	ctx.Redirect(webAppUrl)

	return nil
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
