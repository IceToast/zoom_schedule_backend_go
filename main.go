package main

import (
	// System-Imports
	"log"
	"os"

	// eigene Imports
	"zoom_schedule_backend_go/routes"

	// GitHub Imports
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"github.com/shareed2k/goth_fiber"
	"github.com/subosito/gotenv"
)

const (
	Host = "localhost"
	Port = ":8011"
)

func init() {
	gotenv.Load()
	log.Println(os.Getenv("CONNECTION_STRING"))
}

func main() {
	app := fiber.New()

	goth.UseProviders(
		google.New(os.Getenv("OAUTH_CLIENT_ID"), os.Getenv("OAUTH_SECRET"),
			"http://"+Host+Port+"/api/auth/google/callback"),
		// ... mehr Provider
	)

	app.Use(cors.New())
	api := app.Group("/api") // /api

	// OAuth2-Endpunkte
	api.Get("/auth/:provider", func(ctx *fiber.Ctx) error {
		if gothUser, err := goth_fiber.CompleteUserAuth(ctx); err == nil {
			ctx.JSON(gothUser)
		} else {
			goth_fiber.BeginAuthHandler(ctx)
		}
		return nil
	})
	api.Get("/auth/:provider/callback", func(ctx *fiber.Ctx) error {
		user, err := goth_fiber.CompleteUserAuth(ctx)
		if err != nil {
			log.Fatal(err)
		}

		ctx.JSON(user)
		return nil
	})
	api.Get("/auth/logout/:provider", func(ctx *fiber.Ctx) error {
		if err := goth_fiber.Logout(ctx); err != nil {
			log.Fatal(err)
		}

		ctx.Redirect("/")
		return nil
	})

	// Meeting-Endpunkte
	api.Get("/meeting/:id?", routes.GetMeeting)
	api.Post("/meeting", routes.CreateMeeting)
	api.Put("/meeting/:id", routes.UpdateMeeting)
	api.Delete("/meeting/:id", routes.DeleteMeeting)

	// Test-Endpunkt
	api.Get("/test", func(ctx *fiber.Ctx) error {
		ctx.Format("<p><a href='/api/auth/google'>Google Auth</a></p>")
		return nil
	})

	if err := app.Listen(Port); err != nil {
		log.Fatal(err)
	}
}
