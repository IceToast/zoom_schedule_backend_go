package helpers

import (
	"errors"
	"zoom_schedule_backend_go/db"

	"github.com/gofiber/fiber/v2"
)

func VerifyCookie(ctx *fiber.Ctx) (string, error) {

	store := db.GetStore()
	session, err := store.Get(ctx)
	if err != nil {
		return "", err
	}

	internalUserId, ok := session.Get("internalUserId").(string)

	if !ok || internalUserId == "" {
		session.Destroy()
		return "", errors.New("invalid Cookie or session expired")
	}

	return internalUserId, nil
}
