package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/OVillas/e-commercer-api/config"
	"github.com/OVillas/e-commercer-api/domain"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type sessionRepository struct {
	i           *do.Injector
	db          *gorm.DB
	redisClient *redis.Client
}

func NewSessionRepository(i *do.Injector) (domain.SessionRepository, error) {
	db, err := do.Invoke[*gorm.DB](i)
	if err != nil {
		return nil, err
	}

	redisClient, err := do.Invoke[*redis.Client](i)
	if err != nil {
		return nil, err
	}

	return &sessionRepository{
		i:           i,
		db:          db,
		redisClient: redisClient,
	}, nil
}

func (u *sessionRepository) Create(ctx context.Context, user domain.User, token string) error {
	log := slog.With(
		slog.String("repository", "token"),
		slog.String("func", "Create"),
	)

	log.Info("Initializing token creation process")

	userSession := domain.Session{
		Token:     token,
		Name:      user.Name,
		UserID:    user.ID,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
	}

	userJSON, err := jsoniter.Marshal(userSession)
	if err != nil {
		log.Error("Failed to marshal user data", slog.String("error", err.Error()))
		return err
	}

	if err := u.redisClient.Set(ctx, u.getTokenKey(user.ID.String()), userJSON, time.Duration(config.Env.TokenExp)*time.Hour).Err(); err != nil {
		log.Error("Failed to save token", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (u *sessionRepository) GetUser(ctx context.Context, userID string) (*domain.Session, error) {
	log := slog.With(
		slog.String("repository", "token"),
		slog.String("func", "GetUser"),
	)

	log.Info("Initializing token retrieval process")

	userJSON, err := u.redisClient.Get(ctx, u.getTokenKey(userID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Error("Token not found")
			return nil, domain.ErrSessionNotFound
		}
		log.Error("Failed to retrieve user", slog.String("error", err.Error()))
		return nil, err
	}

	var user domain.Session
	if err := jsoniter.UnmarshalFromString(userJSON, &user); err != nil {
		log.Error("Failed to unmarshal user data", slog.String("error", err.Error()))
		return nil, err
	}

	return &user, nil
}

func (u *sessionRepository) Update(ctx context.Context, user domain.User, token string) error {
	log := slog.With(
		slog.String("repository", "token"),
		slog.String("func", "Update"),
	)

	log.Info("Initializing token update process")

	ttl, err := u.redisClient.TTL(ctx, u.getTokenKey(user.ID.String())).Result()
	if err != nil {
		log.Error("Failed to get TTL for token", slog.String("error", err.Error()))
		return err
	}

	userSession := domain.Session{
		Token:     token,
		Name:      user.Name,
		UserID:    user.ID,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
	}

	userJSON, err := jsoniter.Marshal(userSession)
	if err != nil {
		log.Error("Failed to marshal user data", slog.String("error", err.Error()))
		return err
	}

	if err := u.redisClient.Set(ctx, u.getTokenKey(user.ID.String()), userJSON, ttl).Err(); err != nil {
		log.Error("Failed to save token", slog.String("error", err.Error()))
		return err
	}

	log.Info("Token updated successfully")
	return nil
}

func (u *sessionRepository) SaveOTP(ctx context.Context, email string, otp string) error {
	log := slog.With(
		slog.String("repository", "user_session"),
		slog.String("func", "SaveOTP"),
	)

	log.Info("Initializing OTP save process")

	if err := u.redisClient.Set(ctx, u.getOTPKey(email), otp, time.Duration(config.Env.OTPExp)*time.Minute).Err(); err != nil {
		log.Error("Failed to save OTP", slog.String("error", err.Error()))
		return err
	}

	log.Info("OTP saved successfully")
	return nil
}

func (u *sessionRepository) GetOTP(ctx context.Context, email string) (string, error) {
	log := slog.With(
		slog.String("repository", "user_session"),
		slog.String("func", "VerifyOTP"),
	)

	log.Info("Initializing OTP verification process")

	storedOTP, err := u.redisClient.Get(ctx, u.getOTPKey(email)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Warn("OTP not found")
			return "", domain.ErrOTPNotFound
		}
		log.Error("Failed to retrieve OTP", slog.String("error", err.Error()))
		return "", err
	}

	log.Info("OTP verified successfully")
	return storedOTP, nil
}

func (t *sessionRepository) getTokenKey(id string) string {
	tokenKey := fmt.Sprintf("usersession_%s", id)
	return tokenKey
}

func (t *sessionRepository) getOTPKey(email string) string {
	OTPKey := fmt.Sprintf("usersession_OTP_%s", email)
	return OTPKey
}
