package request

type SendMail struct {
	From     string   `json:"from,omitempty"`
	Subject  string   `json:"subject,omitempty"`
	MailType uint8    `json:"mailType,omitempty"`
	To       []string `json:"to,omitempty"`
}
