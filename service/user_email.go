package service

import (
	"context"
	"html/template"
	"log/slog"
	"os"
	"strings"

	"github.com/OVillas/e-commercer-api/domain"
	"github.com/OVillas/e-commercer-api/secure"
	"github.com/samber/do"
)

type userEmailService struct {
	i                  *do.Injector
	emailService       domain.EmailService
	userSessionService domain.UserSessionService
}

func NewUserEmailService(i *do.Injector) (domain.UserEmailService, error) {
	emailService, err := do.Invoke[domain.EmailService](i)
	if err != nil {
		return nil, err
	}

	userSessionService, err := do.Invoke[domain.UserSessionService](i)
	if err != nil {
		return nil, err
	}

	return &userEmailService{
		i:                  i,
		emailService:       emailService,
		userSessionService: userSessionService,
	}, nil

}

func (u *userEmailService) SendConfirmationCode(ctx context.Context, user domain.User) error {
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

	otp, err := secure.GenerateNumericOTP(key)
	if err != nil {
		log.Error("Failed to generate OTP", slog.String("error", err.Error()))
		return err
	}

	err = u.userSessionService.SaveOTP(ctx, user.Email, otp)
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

	userOTPEmail := domain.UserOTPEmail{
		Name: user.Name,
		OTP:  otp,
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

	_, err = u.emailService.SendEmail(emailReq)
	if err != nil {
		log.Error("Failed to send OTP email", slog.String("error", err.Error()))
		return err
	}

	log.Info("OTP sent successfully")
	return nil

}
