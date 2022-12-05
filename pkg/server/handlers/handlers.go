// Package handlers contains handlers for all Mail API endpoints.
package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/richard-on/auth-service/pkg/authService"
	"github.com/richard-on/mail-service/config"
	"github.com/richard-on/mail-service/internal/templates"
	"github.com/richard-on/mail-service/pkg/logger"
	"github.com/richard-on/mail-service/pkg/server/request"
	"github.com/richard-on/mail-service/pkg/server/response"
	"net/smtp"
)

type MailHandler struct {
	Router      fiber.Router
	AuthService authService.AuthServiceClient
	log         logger.Logger
}

func NewMailHandler(router fiber.Router, authService authService.AuthServiceClient) *MailHandler {
	return &MailHandler{
		Router:      router,
		AuthService: authService,
		log:         logger.NewLogger(config.DefaultWriter, config.LogInfo.Level, "mail-handler"),
	}
}

// Send is endpoint for sending emails
// @Summary      Sends an email
// @Tags         send email
// @Description  Sends an email
// @ID           registration
// @Accept       json
// @Produce      json
// @Param        input        body      request.SendEmail  true  "Send Email request"
// @Success      200          {object}  response.SendSuccess
// @Failure      401,403,500  {object}  response.Error
// @Router       /send [post]
func (h *MailHandler) Send(ctx *fiber.Ctx) error {

	//code, body, err := fasthttp.Post()

	validateRequest := &authService.ValidateRequest{
		AccessToken:  ctx.Cookies("accessToken"),
		RefreshToken: ctx.Cookies("refreshToken"),
	}

	validateResponse, err := h.AuthService.Validate(ctx.Context(), validateRequest)
	if err != nil {
		h.log.Error(err, "validation error")
		return HandleGrpcError(ctx, err)
	}

	sendRequest := &request.SendMail{}

	if err = ctx.BodyParser(sendRequest); err != nil {
		h.log.Debugf("body parsing error: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(response.Error{Error: err.Error()})
	}

	if sendRequest.From != validateResponse.Username {
		h.log.Debug(ErrNoSenderMatch)
		return ctx.Status(fiber.StatusForbidden).JSON(response.Error{Error: ErrNoSenderMatch.Error()})
	}

	//from := "From: " + sendRequest.From + ";\n"
	to := "To: " + sendRequest.To[0] + "\n"
	subject := "Subject: " + sendRequest.Subject + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body, err := templates.SelectTemplate(sendRequest)
	if err != nil {
		h.log.Debugf("template error", err)
		return ctx.Status(fiber.StatusNotImplemented).JSON(response.Error{Error: err.Error()})
	}

	msg := []byte(to + subject + mime + body)

	auth := smtp.PlainAuth("", config.Mailgun.Host, config.Mailgun.Pass, config.SMTP.Host)

	// Sending email
	err = smtp.SendMail(
		config.SMTP.Host+":"+config.SMTP.Port,
		auth,
		validateResponse.Email,
		sendRequest.To,
		msg,
	)
	if err != nil {
		h.log.Error(err, "send error")
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.SendStatus(fiber.StatusOK)
}
