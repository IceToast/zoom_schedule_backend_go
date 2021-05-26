package main

import (
	"log"
	"os"
	"zoom_schedule_backend_go/helpers"
	"zoom_schedule_backend_go/routes"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/shareed2k/goth_fiber"
	"github.com/subosito/gotenv"
)

const (
	Host = "https://zoom.icetoast.cloud"
	Port = ":8011"
)

func init() {
	gotenv.Load()
}

// @title Zoom Schedule Backend
// @version 1.0
// @description Zoom Schedule API using Fiber v2
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host zoom.icetoast.cloud
// @BasePath /
func main() {
	app := fiber.New(fiber.Config{DisableKeepalive: true})

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_SECRET"), Host+"/api/auth/google/callback", "profile", "email"),
		discord.New(os.Getenv("DISCORD_CLIENT_ID"), os.Getenv("DISCORD_SECRET"), Host+"/api/auth/discord/callback", discord.ScopeIdentify, discord.ScopeEmail),
		github.New(os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_SECRET"), Host+"/api/auth/github/callback"),
	)

	app.Use(cors.New(cors.Config(helpers.CorsConfigDefault)))
	app.Static("/docs", "./docs") // Serve static docs/ folder

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

	// User-Endpunkte
	user := api.Group("/user")
	user.Get("/", routes.GetUserData)
	//user.Delete("/", routes.DeleteUser)
	user.Get("/logout", routes.Logout)

	// Meeting-Endpunkte
	meetings := api.Group("/meeting")
	meetings.Get("/", routes.GetMeetings)
	meetings.Post("/", routes.CreateMeeting)
	meetings.Put("/", routes.UpdateMeeting)
	meetings.Delete("/", routes.DeleteMeeting)
	meetings.Delete("/flushSchedule", routes.FlushSchedule)

	// Swagger
	api.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:         Host + "/docs/swagger.json",
		DeepLinking: false,
	}))

	if err := app.Listen(Port); err != nil {
		log.Fatal(err)
	}
}
