package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/richard-on/auth-service/pkg/authService"

	"github.com/richard-on/mail-service/pkg/server/handlers"
)

func MailRouter(app fiber.Router, authClient authService.AuthServiceClient) {

	handler := handlers.NewMailHandler(app, authClient)

	app.Post("/send", handler.Send)
}
