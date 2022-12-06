package request

// SendMail is the request body for sending an email
type SendMail struct {
	From     string      `json:"from,omitempty"`
	Subject  string      `json:"subject,omitempty"`
	To       string      `json:"to,omitempty"`
	Type     string      `json:"templateType,omitempty"`
	Template interface{} `json:"template,omitempty"`
}
