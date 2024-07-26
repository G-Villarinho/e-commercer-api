package service

import (
	"log/slog"

	"github.com/OVillas/e-commercer-api/domain"
	"github.com/resend/resend-go/v2"
	"github.com/samber/do"
)

type emailService struct {
	i      *do.Injector
	client *resend.Client
}

func NewEmailService(i *do.Injector) (domain.EmailService, error) {
	client, err := do.Invoke[*resend.Client](i)
	if err != nil {
		return nil, err
	}

	return &emailService{
		i:      i,
		client: client,
	}, nil
}

func (e *emailService) SendEmail(request domain.SendEmailRequest) (*domain.SendEmailResponse, error) {
	log := slog.With(
		slog.String("service", "email"),
		slog.String("func", "SendEmail"),
	)

	log.Info("Starting to send email")

	params := &resend.SendEmailRequest{
		From:    request.From,
		To:      request.To,
		Html:    request.Html,
		Subject: request.Subject,
		Cc:      request.Cc,
		Bcc:     request.Bcc,
		ReplyTo: request.ReplyTo,
	}

	sent, err := e.client.Emails.Send(params)
	if err != nil {
		log.Error("Fail to send email", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("Email sent successfully", slog.String("emailID", sent.Id))
	return &domain.SendEmailResponse{Id: sent.Id}, nil
}
