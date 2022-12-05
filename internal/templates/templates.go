package templates

import (
	"errors"
	"github.com/matcornic/hermes/v2"
	"github.com/richard-on/mail-service/internal/model"
	"github.com/richard-on/mail-service/pkg/server/request"
	"os"
)

func SelectTemplate(req *request.SendMail) (string, error) {
	switch req.MailType {
	case model.TaskCoordination:
		return CoordinationTemplate(req)
	default:
		return "", errors.New("not implemented")

	}
}

func CoordinationTemplate(req *request.SendMail) (string, error) {
	h := hermes.Hermes{
		// Optional Theme
		// Theme: new(Default)
		Product: hermes.Product{
			// Appears in header & footer of e-mails
			Name: "Coordination Service",
			Link: "https://richardhere.dev/",
			Logo: "https://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
		},
	}

	email := hermes.Email{
		Body: hermes.Body{
			Name:   req.From,
			Intros: []string{req.Subject},
			Actions: []hermes.Action{
				{
					Instructions: "To approve this task:",
					Button: hermes.Button{
						Color: "#22BC66", // Optional action button color
						Text:  "Approve",
						Link:  "https://richardhere.dev/tags/go",
					},
				},
				{
					Instructions: "To decline this task:",
					Button: hermes.Button{
						Color: "#BC3922", // Optional action button color
						Text:  "Decline",
						Link:  "https://richardhere.dev/tags/hugo",
					},
				},
			},
		},
	}

	emailBody, err := h.GenerateHTML(email)
	if err != nil {
		return "", err
	}

	// Optionally, preview the generated HTML e-mail by writing it to a local file
	err = os.WriteFile("request.html", []byte(emailBody), 0644)
	if err != nil {
		return "", err
	}

	return emailBody, nil
}
