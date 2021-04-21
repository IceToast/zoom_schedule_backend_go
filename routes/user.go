package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"zoom_schedule_backend_go/db"
	"zoom_schedule_backend_go/helpers"

	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

const collectionSession = "session"

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

func DeleteUser(ctx *fiber.Ctx) error {
	//Verify Cookie
	//internalUserId, err := helpers.VerifyCookie(ctx)
	//if err != nil {
	//	return ctx.Status(403).SendString(err.Error())
	//}

	collectionSession, err := db.GetMongoDbCollection(dbName, collectionSession)
	if err != nil {
		return ctx.SendStatus(500)
	}

	sessions, err := collectionSession.Find(context.TODO(), bson.M{})
	if err != nil {
		return ctx.SendStatus(500)
	}

	bool := sessions.Next(context.TODO())
	fmt.Println(bool, sessions)

	fmt.Println("here1")
	for !sessions.Next(context.TODO()) {
		fmt.Println("here2")

		//var session struct {
		//	Id    primitive.ObjectID `json:"_id" bson:"_id"`
		//	Value primitive.Binary   `json:"value" bson:"value"`
		//}

		var cursor bson.M
		fmt.Println("here2")

		err := sessions.Decode(&cursor)
		if err != nil {
			log.Default()
		}
		fmt.Println(cursor)
	}
	defer sessions.Close(context.TODO())

	//deleteInternalUserErr := DeleteInternalUser(internalUserId)
	//if deleteInternalUserErr != nil {
	//	return ctx.SendStatus(500)
	//}
	//deleteExternalUserErr := DeleteExternalUser(internalUserId)
	//if deleteExternalUserErr != nil {
	//	return ctx.SendStatus(500)
	//}

	return ctx.SendStatus(200)

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

func DeleteSession(sessionId primitive.ObjectID) error {

	return nil
}

// Logout godoc
// @Summary Deletes user session, expires session_id cookie and redirects to startpage
// @Description Resolves a userId via a given session cookie.
// @Produce json
// @Success 200
// @Failure 403 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /api/user/logout [get]
func Logout(ctx *fiber.Ctx) error {
	//Verify Cookie
	_, err := helpers.VerifyCookie(ctx)
	if err != nil {
		return ctx.Status(403).SendString(err.Error())
	}

	store := db.GetStore()
	session, err := store.Get(ctx)
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}

	session.Destroy()
	ctx.Redirect(webAppUrl)
	return nil
}

// GetUserData godoc
// @Summary Retrieves user data of logged in user
// @Description Resolves a userId via a given session cookie.
// @Produce json
// @Success 200 {object} userData
// @Failure 403 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /api/user [get]
func GetUserData(ctx *fiber.Ctx) error {
	//Verify Cookie
	internalUserId, err := helpers.VerifyCookie(ctx)
	if err != nil {
		return ctx.Status(403).SendString(err.Error())
	}

	userCollection, err := db.GetMongoDbCollection(dbName, collectionUser)
	if err != nil {
		return ctx.SendStatus(500)
	}

	externalAuthCollection, err := db.GetMongoDbCollection(dbName, collectionExternalAuth)
	if err != nil {
		return ctx.SendStatus(500)
	}

	internalUserObjId, _ := primitive.ObjectIDFromHex(internalUserId)

	var internalUser *User
	err = userCollection.FindOne(context.Background(), bson.M{"_id": internalUserObjId}).Decode(&internalUser)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return ctx.Status(500).SendString("User not found")
		}
	}
	fmt.Println("here1", internalUser)

	var externalAuthUser *ExternalAuthUser
	err = externalAuthCollection.FindOne(context.Background(), bson.M{"internaluserid": internalUserId}).Decode(&externalAuthUser)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return ctx.Status(500).SendString("User not found")
		}
	}
	fmt.Println("here2", externalAuthUser)

	userdata := userData{
		Id:        internalUserId,
		Username:  internalUser.UserName,
		Email:     internalUser.Email,
		AvatarUrl: externalAuthUser.AvatarURL,
		Platform:  externalAuthUser.Platform,
	}

	response, _ := json.Marshal(userdata)
	ctx.Send(response)
	return nil

}
