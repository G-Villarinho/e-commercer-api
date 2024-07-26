package service

import (
	"context"
	"html/template"
	"log/slog"
	"os"
	"strings"

	"github.com/OVillas/e-commercer-api/domain"
	"github.com/OVillas/e-commercer-api/secure"
	"github.com/resend/resend-go/v2"
	"github.com/samber/do"
)

type emailService struct {
	i                  *do.Injector
	client             *resend.Client
	userSessionService domain.SessionService
}

func NewEmailService(i *do.Injector) (domain.EmailService, error) {
	client, err := do.Invoke[*resend.Client](i)
	if err != nil {
		return nil, err
	}

	userSessionService, err := do.Invoke[domain.SessionService](i)
	if err != nil {
		return nil, err
	}

	return &emailService{
		i:                  i,
		client:             client,
		userSessionService: userSessionService,
	}, nil
}

func (e *emailService) SendConfirmationCode(ctx context.Context, user domain.User) error {
	log := slog.With(
		slog.String("service", "userEmail"),
		slog.String("func", "SendConfirmationCode"),
	)

	log.Info("Initializing send connfirmation code in process")

	key, err := secure.GenerateSecret(user.Email)
	if err != nil {
		log.Error("Failed to generate key user", slog.String("error", err.Error()))
		return err
	}

	OTP, err := secure.GenerateNumericOTP(key)
	if err != nil {
		log.Error("Failed to generate OTP", slog.String("error", err.Error()))
		return err
	}

	err = e.userSessionService.SaveOTP(ctx, user.Email, OTP)
	if err != nil {
		log.Error("Failed to save OTP", slog.String("error", err.Error()))
		return err
	}

	tmpl, err := os.ReadFile("templates\\otp_template.html")
	if err != nil {
		log.Error("Failed to read email template", slog.String("error", err.Error()))
		return err
	}

	t, err := template.New("emailTemplate").Parse(string(tmpl))
	if err != nil {
		log.Error("Failed to parse email template", slog.String("error", err.Error()))
		return err
	}

	userOTPEmail := domain.OTPEmailPayload{
		Name: user.Name,
		OTP:  OTP,
	}

	var body strings.Builder
	err = t.Execute(&body, userOTPEmail)
	if err != nil {
		log.Error("Failed to execute email template", slog.String("error", err.Error()))
		return err
	}

	emailReq := domain.SendEmailRequest{
		From:    "Acme <onboarding@resend.dev>",
		To:      []string{user.Email},
		Subject: "Your OTP Code",
		Html:    body.String(),
	}

	_, err = e.sendEmail(emailReq)
	if err != nil {
		log.Error("Failed to send OTP email", slog.String("error", err.Error()))
		return err
	}

	log.Info("OTP sent successfully")
	return nil

}

func (e *emailService) sendEmail(request domain.SendEmailRequest) (*domain.SendEmailResponse, error) {
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
