package templates

import (
	"github.com/matcornic/hermes/v2"
	"github.com/richard-on/mail-service/pkg/server/request"
)

// Info is the template for info emails
type Info struct {
	Body string `json:"body,omitempty"`
}

func (i *Info) setTemplate(req *request.SendMail) (string, string, error) {
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
			Intros: []string{i.Body},
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
