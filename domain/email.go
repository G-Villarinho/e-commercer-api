package domain

import "context"

type SendEmailRequest struct {
	From    string
	To      []string
	Html    string
	Subject string
	Cc      []string
	Bcc     []string
	ReplyTo string
}

type OTPEmailPayload struct {
	Name string
	OTP  string
}

type SendEmailResponse struct {
	Id string
}

type EmailService interface {
	SendConfirmationCode(ctx context.Context, user User) error
}
