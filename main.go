package main

import (
	"zoom_schedule_backend_go/routes"

	"github.com/gofiber/fiber"
	"github.com/subosito/gotenv"
)

const port = 8011
const dbName = "zoom_schedule"
const collectionMeeting = "meeting"

func init() {
	gotenv.Load()
}

func main() {
	app := fiber.New()
	app.Get("/meeting/:id?", routes.GetMeeting)
	app.Post("/meeting", routes.CreateMeeting)
	app.Put("/meeting/:id", routes.UpdateMeeting)
	app.Delete("/meeting/:id", routes.DeleteMeeting)
	app.Listen(port)
  }