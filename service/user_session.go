package service

import (
	"context"
	"log/slog"

	"github.com/OVillas/e-commercer-api/domain"
	"github.com/OVillas/e-commercer-api/middleware"
	"github.com/OVillas/e-commercer-api/util"
	"github.com/samber/do"
)

type userSessionService struct {
	i                     *do.Injector
	userSessionRepository domain.UserSessionRepository
	userRepository        domain.UserRepository
}

func NewUserSessionService(i *do.Injector) (domain.UserSessionService, error) {
	tokenRepository, err := do.Invoke[domain.UserSessionRepository](i)
	if err != nil {
		return nil, err
	}

	userRepository, err := do.Invoke[domain.UserRepository](i)
	if err != nil {
		return nil, err
	}

	return &userSessionService{
		i:                     i,
		userSessionRepository: tokenRepository,
		userRepository:        userRepository,
	}, nil
}

func (t *userSessionService) Create(ctx context.Context, user domain.User) (string, error) {
	log := slog.With(
		slog.String("service", "token"),
		slog.String("func", "Create"),
	)

	log.Info("Initializing token creation process")

	token, err := util.CreateToken(user)
	if err != nil {
		log.Error("Failed to create token", slog.String("error", err.Error()))
		return "", err
	}

	if err := t.userSessionRepository.Create(ctx, user, token); err != nil {
		log.Error("Failed to save token", slog.String("error", err.Error()))
		return "", err
	}

	return token, nil
}

func (t *userSessionService) GetUser(ctx context.Context, token string) (*domain.UserSession, error) {
	log := slog.With(
		slog.String("service", "token"),
		slog.String("func", "GetUser"),
	)

	log.Info("Initializing token retrieval process")

	userID, err := util.ExtractUserIDFromToken(token)
	if err != nil {
		log.Error("Failed to extract user ID from token", slog.String("error", err.Error()))
		return nil, err
	}

	user, err := t.userSessionRepository.GetUser(ctx, userID)
	if err != nil {
		log.Error("Failed to retrieve user", slog.String("error", err.Error()))
		return nil, err
	}

	if user == nil {
		log.Error("Session not found")
		return nil, domain.ErrSessionNotFound
	}

	if token != user.Token {
		log.Error("Session mismatch")
		return nil, domain.ErrTokenInvalid
	}

	return user, nil
}

func (t *userSessionService) Update(ctx context.Context) error {
	log := slog.With(
		slog.String("service", "token"),
		slog.String("func", "Update"),
	)

	log.Info("Initializing token update process")

	userSession, ok := ctx.Value(middleware.UserKey).(*domain.UserSession)
	if !ok || userSession == nil {
		return domain.ErrUserNotFoundInContext
	}

	user, err := t.userRepository.GetByID(ctx, userSession.UserID)
	if err != nil {
		log.Error("Failed to get user by ID", slog.String("error", err.Error()))
		return err
	}

	if err := t.userSessionRepository.Update(ctx, *user, userSession.Token); err != nil {
		log.Error("Failed to update token", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (t *userSessionService) SaveOTP(ctx context.Context, email string, otp string) error {
	log := slog.With(
		slog.String("service", "token"),
		slog.String("func", "SaveOTP"),
	)

	log.Info("Initializing OTP save process")

	err := t.userSessionRepository.SaveOTP(ctx, email, otp)
	if err != nil {
		log.Error("Failed to save OTP", slog.String("error", err.Error()))
		return err
	}

	log.Info("OTP saved successfully")
	return nil
}
