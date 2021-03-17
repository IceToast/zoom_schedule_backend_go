package main

import (
	"zoom_schedule_backend_go/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/subosito/gotenv"
)

const port = ":8011"

func init() {
	gotenv.Load()
}

func main() {

	app := fiber.New()
	app.Use(cors.New())
	app.Get("/meeting/:id?", routes.GetMeeting)
	app.Post("/meeting", routes.CreateMeeting)
	app.Put("/meeting/:id", routes.UpdateMeeting)
	app.Delete("/meeting/:id", routes.DeleteMeeting)
	app.Listen(port)
  }