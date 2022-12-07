package templates

import (
	"github.com/matcornic/hermes/v2"
	"github.com/richard-on/mail-service/pkg/server/request"
)

// Coordination is the template for coordination emails
type Coordination struct {
	acceptLink  string `json:"acceptLink"`
	declineLink string `json:"declineLink"`
}

func (c *Coordination) setTemplate(req *request.SendMail) (string, string, error) {
	h := hermes.Hermes{
		Product: hermes.Product{
			Name: "Coordination Service",
			Link: "https://richardhere.dev/",
			Logo: "https://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
		},
	}

	email := hermes.Email{
		Body: hermes.Body{
			Name:   req.From,
			Intros: []string{"This is coordination service"},
			Actions: []hermes.Action{
				{
					Instructions: "To approve this task:",
					Button: hermes.Button{
						Color: "#22BC66", // Optional action button color
						Text:  "Approve",
						Link:  c.acceptLink,
					},
				},
				{
					Instructions: "To decline this task:",
					Button: hermes.Button{
						Color: "#BC3922", // Optional action button color
						Text:  "Decline",
						Link:  c.declineLink,
					},
				},
			},
		},
	}

	emailHTML, err := h.GenerateHTML(email)
	if err != nil {
		return "", "", err
	}

	emailPlain, err := h.GeneratePlainText(email)
	if err != nil {
		return "", "", err
	}

	return emailHTML, emailPlain, nil
}
