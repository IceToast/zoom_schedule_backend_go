package routes

//import (
//	"context"
//	"encoding/json"
//
//	"zoom_schedule_backend_go/db"
//
//	"github.com/gofiber/fiber/v2"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/bson/primitive"
//)
//
//type User struct {
//	_id      primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
//	UserName     string             `json:"name,omitempty" bson":name.omitempty"`
//	Email     string             `json:"link,omitempty" bson":link.omitempty`
//	Days []Day             `json:"days,omitempty" bson:"days,omitempty"`
//}
//
//type Meeting struct {
//	_id      primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
//	Name     string             `json:"name,omitempty" bson":name.omitempty"`
//	Link     string             `json:"link,omitempty" bson":link.omitempty`
//	Password string             `json:"password,omitempty" bson:"password,omitempty"`
//}
//
//type Day struct {
//	_id primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
//	Name string `json:"name,omitempty" bson:"name.omitempty"`
//	Meetings []Meeting `json:"meetings,omitempty" bson:"meetings.omitempty"`
//}
//
//const dbName = "zoom_schedule"
//const collectionName = "user"
