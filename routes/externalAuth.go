package routes

import (
	"zoom_schedule_backend_go/db"

	"github.com/gofiber/fiber/v2"
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

// OAuth godoc
// @Summary This is not a request Route! You must redirect to this route.
// @Description Redirects to OAuth provider to start Auth - leads to OAuth callback and sets session cookie - Currently Supported: Discord, Google, Github
// @Router /api/auth/{provider} [get]
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

	return ctx.Redirect(webAppUrl)
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
	db.CloseMongoDbConnection(collection)

	return result, nil
}

func DeleteExternalUser(internalUserId string) error {
	collection, err := db.GetMongoDbCollection(dbName, collectionExternalAuth)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"internaluserid": internalUserId})
	if err != nil {
		return err
	}
	db.CloseMongoDbConnection(collection)

	return nil
}
