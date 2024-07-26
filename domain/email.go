package domain

type SendEmailRequest struct {
	From    string
	To      []string
	Html    string
	Subject string
	Cc      []string
	Bcc     []string
	ReplyTo string
}

type SendEmailResponse struct {
	Id string
}

type EmailService interface {
	SendEmail(req SendEmailRequest) (*SendEmailResponse, error)
}
