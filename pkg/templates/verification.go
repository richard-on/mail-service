package templates

import (
	"github.com/matcornic/hermes/v2"
	"github.com/richard-on/mail-service/pkg/server/request"
)

// Verification is the template for verification emails
type Verification struct {
	VerifyLink string `json:"verifyLink"`
}

func (v *Verification) setTemplate(req *request.SendMail) (string, string, error) {
	h := hermes.Hermes{
		Product: hermes.Product{
			Name: "Info",
			Link: "https://richardhere.dev/",
			Logo: "https://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
		},
	}

	email := hermes.Email{
		Body: hermes.Body{
			Name:   req.From,
			Intros: []string{"Verification"},
			Actions: []hermes.Action{
				{
					Instructions: "To approve this task:",
					Button: hermes.Button{
						Color: "#22BC66", // Optional action button color
						Text:  "Approve",
						Link:  v.VerifyLink,
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
