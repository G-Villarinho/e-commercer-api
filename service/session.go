package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/OVillas/e-commercer-api/config"
	"github.com/OVillas/e-commercer-api/domain"
	"github.com/OVillas/e-commercer-api/middleware"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
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

	token, err := t.createToken(user)
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

	sessionToken, err := t.extractSessionFromToken(token)
	if err != nil {
		log.Error("Failed to extract user ID from token", slog.String("error", err.Error()))
		return nil, err
	}

	session, err := t.sessionRepository.GetUser(ctx, sessionToken.UserID.String())
	if err != nil {
		log.Error("Failed to retrieve user", slog.String("error", err.Error()))
		return nil, err
	}

	if session == nil {
		log.Error("Session not found")
		return nil, domain.ErrSessionNotFound
	}

	if token != session.Token {
		log.Error("Session mismatch")
		return nil, domain.ErrTokenInvalid
	}

	return session, nil
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

func (t *sessionService) GetOTP(ctx context.Context, email string) (string, error) {
	log := slog.With(
		slog.String("service", "token"),
		slog.String("func", "getOTP"),
	)

	log.Info("Initializing OTP get process")
	OTP, err := t.sessionRepository.GetOTP(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrOTPNotFound) {
			log.Warn("otp not found")
			return "", nil
		}

		return "", err
	}

	log.Info("Initializing OTP get completed")
	return OTP, nil
}

func (t *sessionService) createToken(user domain.User) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"id":        user.ID,
		"name":      user.Name,
		"email":     user.Email,
		"avatarURL": user.AvatarURL,
	})

	tokenString, err := token.SignedString(config.Env.PrivateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t *sessionService) extractSessionFromToken(tokenString string) (*domain.Session, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, domain.ErrTokenInvalid
	}

	jsonStr, err := jsoniter.Marshal(claims)
	if err != nil {
		return nil, err
	}

	userIDStr, ok := claims["id"].(string)
	if !ok {
		return nil, domain.ErrTokenInvalid
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, domain.ErrTokenInvalid
	}

	var session domain.Session
	err = jsoniter.Unmarshal(jsonStr, &session)
	if err != nil {
		return nil, err
	}
	session.UserID = userID

	return &session, nil
}
