package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const dbName = "zoom_schedule"
const collectionSessions = "sessions"

//GetMongoDbConnection get connection of mongodb
func GetMongoDbConnection() (*mongo.Client, error) {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("CONNECTION_STRING")))

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, nil
}

func GetMongoDbCollection(DbName string, CollectionName string) (*mongo.Collection, error) {
	client, err := GetMongoDbConnection()

	if err != nil {
		return nil, err
	}

	collection := client.Database(DbName).Collection(CollectionName)
	return collection, nil
}

func CloseMongoDbConnection(collection *mongo.Collection) {
	collection.Database().Client().Disconnect(context.TODO())
}

func GetStore() *session.Store {
	// Fiber Middleware Storage
	storage := mongodb.New(mongodb.Config{
		ConnectionURI: os.Getenv("CONNECTION_STRING"),
		Database:      dbName,
		Collection:    collectionSessions,
		Reset:         false,
	})

	// Storage: MongoDB, Expiration: 30days
	store := session.New(session.Config{
		Storage:        storage,
		Expiration:     720 * time.Hour,
		CookieSecure:   true,
		CookieSameSite: "None",
	})

	return store
}
