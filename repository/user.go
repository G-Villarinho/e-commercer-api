package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/GSVillas/e-commercer-api/domain"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type userRepository struct {
	i           *do.Injector
	db          *gorm.DB
	redisClient *redis.Client
}

func NewUserRepository(i *do.Injector) (domain.UserRepository, error) {
	db, err := do.Invoke[*gorm.DB](i)
	if err != nil {
		return nil, err
	}

	redisClient, err := do.Invoke[*redis.Client](i)
	if err != nil {
		return nil, err
	}

	return &userRepository{
		i:           i,
		db:          db,
		redisClient: redisClient,
	}, nil
}

func (u *userRepository) Create(ctx context.Context, user domain.User) error {
	log := slog.With(
		slog.String("repository", "user"),
		slog.String("func", "Create"),
	)

	log.Info("Initializing user creation process")

	if err := u.db.WithContext(ctx).Create(&user).Error; err != nil {
		log.Error("Failed to create user", slog.String("error", err.Error()))
		return err
	}

	log.Info("User created successfully")
	return nil
}

func (u *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	log := slog.With(
		slog.String("repository", "user"),
		slog.String("func", "GetByEmail"),
	)

	log.Info("Initializing get user by email process")

	var user domain.User
	if err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("User not found")
			return nil, nil
		}

		log.Error("Failed to get user by email", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("User found successfully")
	return &user, nil
}

func (u *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	log := slog.With(
		slog.String("repository", "user"),
		slog.String("func", "GetByID"),
	)

	log.Info("Initializing get user by ID process")

	var user domain.User
	if err := u.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("User not found")
			return nil, nil
		}

		log.Error("Failed to get user by ID", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("User found successfully")
	return &user, nil
}

func (u *userRepository) UpdateName(ctx context.Context, id uuid.UUID, name string) error {
	log := slog.With(
		slog.String("repository", "user"),
		slog.String("func", "UpdateName"),
	)

	log.Info("Initializing user name update process")

	if err := u.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Update("name", name).Error; err != nil {
		log.Error("Failed to update user name", slog.String("error", err.Error()))
		return err
	}

	log.Info("User name updated successfully")
	return nil
}

func (u *userRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	log := slog.With(
		slog.String("repository", "user"),
		slog.String("func", "UpdatePassword"),
	)

	log.Info("Initializing user password update process")

	if err := u.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Update("passwordHash", passwordHash).Error; err != nil {
		log.Error("Failed to update user password", slog.String("error", err.Error()))
		return err
	}

	log.Info("User password updated successfully")
	return nil
}

func (u *userRepository) UpdateConfirmEmail(ctx context.Context, id uuid.UUID) error {
	log := slog.With(
		slog.String("repository", "user"),
		slog.String("func", "UpdateConfirmEmail"),
	)

	log.Info("Initializing user name update process")

	if err := u.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Update("emailConfirmed", true).Error; err != nil {
		log.Error("Failed to update user name", slog.String("error", err.Error()))
		return err
	}

	log.Info("User email confirm updated successfully")
	return nil
}
