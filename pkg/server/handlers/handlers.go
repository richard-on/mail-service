// Package handlers contains handlers for all Mail API endpoints.
package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/richard-on/auth-service/pkg/authService"
	"github.com/richard-on/mail-service/config"
	"github.com/richard-on/mail-service/pkg/logger"
	"github.com/richard-on/mail-service/pkg/server/request"
	"github.com/richard-on/mail-service/pkg/server/response"
	templates2 "github.com/richard-on/mail-service/pkg/templates"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"os"
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

	html, plain, err := templates2.GetTemplate(sendRequest)
	if err == templates2.ErrNoSuchTemplate {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.Error{Error: templates2.ErrNoSuchTemplate.Error()})
	} else if err == templates2.ErrBadFormat || errors.Unwrap(err) == templates2.ErrBadFormat {
		h.log.Debug(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(response.Error{Error: templates2.ErrBadFormat.Error()})
	} else if err != nil {
		h.log.Debugf("template error", err)
		return ctx.Status(fiber.StatusNotImplemented).JSON(response.Error{Error: err.Error()})
	}

	from := mail.NewEmail(sendRequest.From, "no-reply@richardhere.dev")
	subject := sendRequest.Subject
	to := mail.NewEmail(sendRequest.To, sendRequest.To)
	message := mail.NewSingleEmail(from, subject, to, plain, html)
	client := sendgrid.NewSendClient(os.Getenv("MAILGUN_PASS"))
	sendResponse, err := client.Send(message)
	if err != nil || sendResponse.StatusCode != 202 {
		h.log.Error(err, "send error")
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	h.log.Debug(sendResponse)

	return ctx.SendStatus(fiber.StatusOK)
}
