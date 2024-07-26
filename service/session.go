package service

import (
	"context"
	"log/slog"

	"github.com/OVillas/e-commercer-api/domain"
	"github.com/OVillas/e-commercer-api/middleware"
	"github.com/OVillas/e-commercer-api/util"
	"github.com/samber/do"
)

type sessionService struct {
	i                 *do.Injector
	sessionRepository domain.SessionRepository
	userRepository    domain.UserRepository
}

func NewSessionService(i *do.Injector) (domain.SessionService, error) {
	sessionRepository, err := do.Invoke[domain.SessionRepository](i)
	if err != nil {
		return nil, err
	}

	userRepository, err := do.Invoke[domain.UserRepository](i)
	if err != nil {
		return nil, err
	}

	return &sessionService{
		i:                 i,
		sessionRepository: sessionRepository,
		userRepository:    userRepository,
	}, nil
}

func (t *sessionService) Create(ctx context.Context, user domain.User) (string, error) {
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

	if err := t.sessionRepository.Create(ctx, user, token); err != nil {
		log.Error("Failed to save token", slog.String("error", err.Error()))
		return "", err
	}

	return token, nil
}

func (t *sessionService) GetUser(ctx context.Context, token string) (*domain.Session, error) {
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

	user, err := t.sessionRepository.GetUser(ctx, userID)
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

func (t *sessionService) Update(ctx context.Context) error {
	log := slog.With(
		slog.String("service", "token"),
		slog.String("func", "Update"),
	)

	log.Info("Initializing token update process")

	userSession, ok := ctx.Value(middleware.UserKey).(*domain.Session)
	if !ok || userSession == nil {
		return domain.ErrUserNotFoundInContext
	}

	user, err := t.userRepository.GetByID(ctx, userSession.UserID)
	if err != nil {
		log.Error("Failed to get user by ID", slog.String("error", err.Error()))
		return err
	}

	if err := t.sessionRepository.Update(ctx, *user, userSession.Token); err != nil {
		log.Error("Failed to update token", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (t *sessionService) SaveOTP(ctx context.Context, email string, otp string) error {
	log := slog.With(
		slog.String("service", "token"),
		slog.String("func", "SaveOTP"),
	)

	log.Info("Initializing OTP save process")

	err := t.sessionRepository.SaveOTP(ctx, email, otp)
	if err != nil {
		log.Error("Failed to save OTP", slog.String("error", err.Error()))
		return err
	}

	log.Info("OTP saved successfully")
	return nil
}
