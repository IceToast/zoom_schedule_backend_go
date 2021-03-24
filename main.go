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
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/google"
	"github.com/shareed2k/goth_fiber"
	"github.com/subosito/gotenv"
)

const (
	Host = "zoomapi.icetoast.cloud"
	Port = ":8011"
)

func init() {
	gotenv.Load()
}

func main() {
	app := fiber.New()

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_SECRET"), "https://"+Host+"/api/auth/google/callback"),
		discord.New(os.Getenv("DISCORD_CLIENT_ID"), os.Getenv("DISCORD_SECRET"), "https://"+Host+"/api/auth/discord/callback", discord.ScopeIdentify, discord.ScopeEmail),
		//The Github method at the time is depricated, wait for next Goth-Release
		//github.New(os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_SECRET"), "https://"+Host+"/api/auth/github/callback"),
	)

	app.Use(cors.New())
	api := app.Group("/api")

	// OAuth2-Endpunkte
	api.Get("/auth/:provider", goth_fiber.BeginAuthHandler)
	api.Get("/auth/:provider/callback", func(ctx *fiber.Ctx) error {
		user, err := goth_fiber.CompleteUserAuth(ctx)
		if err != nil {
			return ctx.SendString(err.Error())
		}

		ctx.JSON(user)
		return nil
	})
	api.Get("/auth/logout/:provider", func(ctx *fiber.Ctx) error {
		if err := goth_fiber.Logout(ctx); err != nil {
			return ctx.SendString(err.Error())
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
