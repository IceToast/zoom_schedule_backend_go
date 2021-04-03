package main

import (
	"log"
	"os"
	"zoom_schedule_backend_go/routes"

	"github.com/arsmn/fiber-swagger/v2"
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

// @title Zoom Schedule Backend
// @version 1.0
// @description The backend based on Fiber
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host zoomapi.icetoast.cloud:8011
// @BasePath /api
func main() {
	app := fiber.New()

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_SECRET"), "https://"+Host+"/api/auth/google/callback", "profile", "email"),
		discord.New(os.Getenv("DISCORD_CLIENT_ID"), os.Getenv("DISCORD_SECRET"), "https://"+Host+"/api/auth/discord/callback", discord.ScopeIdentify, discord.ScopeEmail),
		//The Github method at the time is deprecated, wait for next Goth-Release
		//github.New(os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_SECRET"), "https://"+Host+"/api/auth/github/callback"),
	)

	app.Use(cors.New())

	api := app.Group("/api")

	// OAuth2-Endpunkte
	auth := api.Group("/auth")
	auth.Get("/:provider", goth_fiber.BeginAuthHandler)
	auth.Get("/:provider/callback", routes.ProviderCallback)
	auth.Get("/logout/:provider", func(ctx *fiber.Ctx) error {
		if err := goth_fiber.Logout(ctx); err != nil {
			return ctx.SendString(err.Error())
		}

		ctx.Redirect("/")
		return nil
	})

	// Meeting-Endpunkte
	meetings := api.Group("/meeting")
	meetings.Get("/", routes.GetMeeting)
	meetings.Post("/", routes.CreateMeeting)
	meetings.Put("/", routes.UpdateMeeting)
	meetings.Delete("/", routes.DeleteMeeting)

	// Swagger
	app.Get("/swagger/*", swagger.Handler)

	if err := app.Listen(Port); err != nil {
		log.Fatal(err)
	}
}
