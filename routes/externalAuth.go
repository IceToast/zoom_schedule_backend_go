package routes

import (
	"errors"
	"net/http"
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
)

// OAuth godoc
// @Summary This is not a request Route! You must redirect to this route.
// @Description Redirects to OAuth provider to start Auth
// @Description Redirects to OAuth callback and sets session cookie
// @Description Providers currently Supported: discord, google, github
// @Param redirectUrl path string true "Url api redirects after authentication"
// @Router /api/auth/{provider} [get]
func ProviderCallback(ctx *fiber.Ctx) error {
	redirectUrl, err := goth_fiber.GetFromSession("redirectUrl", ctx)
	if err != nil {
		return ctx.SendString(err.Error())
	}
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

	GetSession(ctx, externalUser, redirectUrl)

	return ctx.Redirect(redirectUrl)
}

func GetExternalUser(externaluserID string) (*ExternalAuthUser, error) {
	collection := db.DbInstance.DB.Collection(collectionExternalAuth)

	var result *ExternalAuthUser
	err := collection.FindOne(context.TODO(), bson.M{"externaluserid": externaluserID}).Decode(&result)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
	}

	return result, nil
}

func DeleteExternalUser(internalUserId string) error {
	collection := db.DbInstance.DB.Collection(collectionExternalAuth)

	_, err := collection.DeleteOne(context.Background(), bson.M{"internaluserid": internalUserId})
	if err != nil {
		return err
	}

	return nil
}

func BeginAuthHandler(ctx *fiber.Ctx) error {
	url, err := GetAuthURL(ctx)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return err
	}

	return ctx.Redirect(url, http.StatusTemporaryRedirect)
}

func GetAuthURL(ctx *fiber.Ctx) (string, error) {
	if goth_fiber.SessionStore == nil {
		return "", goth_fiber.ErrSessionNil
	}

	providerName, err := goth_fiber.GetProviderName(ctx)
	if err != nil {
		return "", err
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return "", err
	}

	redirectUrl, err := GetRedirectUrl(ctx)
	if err != nil {
		return "", err
	}

	sess, err := provider.BeginAuth(goth_fiber.SetState(ctx))
	if err != nil {
		return "", err
	}

	url, err := sess.GetAuthURL()
	if err != nil {
		return "", err
	}

	err = goth_fiber.StoreInSession(providerName, sess.Marshal(), ctx)
	if err != nil {
		return "", err
	}

	err = goth_fiber.StoreInSession("redirectUrl", redirectUrl, ctx)
	if err != nil {
		return "", err
	}

	return url, err
}

func GetRedirectUrl(ctx *fiber.Ctx) (string, error) {
	// try to get it from the url param "redirectUrl"
	if p := ctx.Query("redirectUrl"); p != "" {
		return p, nil
	}
	// try to get it from the url param ":redirectUrl"
	if p := ctx.Params("redirectUrl"); p != "" {
		return p, nil
	}

	return "", errors.New("you must provide a redirectUrl")
}
