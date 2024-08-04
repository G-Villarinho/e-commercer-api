package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/OVillas/e-commercer-api/domain"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type storeRepository struct {
	i           *do.Injector
	db          *gorm.DB
	redisClient *redis.Client
}

func NewStoreRepository(i *do.Injector) (domain.StoreRepository, error) {
	db, err := do.Invoke[*gorm.DB](i)
	if err != nil {
		return nil, err
	}

	redisClient, err := do.Invoke[*redis.Client](i)
	if err != nil {
		return nil, err
	}

	return &storeRepository{
		i:           i,
		db:          db,
		redisClient: redisClient,
	}, nil
}

func (s *storeRepository) Create(ctx context.Context, store domain.Store) error {
	log := slog.With(
		slog.String("func", "Create"),
		slog.String("repository", "store"),
	)

	log.Info("Initializing store creation process")

	if err := s.db.WithContext(ctx).Create(&store).Error; err != nil {
		log.Error("Failed to create store", slog.String("error", err.Error()))
		return err
	}

	log.Info("store creation excuted sucessfully")
	return nil
}

func (s *storeRepository) GetAll(ctx context.Context, userID uuid.UUID) ([]*domain.Store, error) {
	log := slog.With(
		slog.String("repository", "store"),
		slog.String("func", "GetAll"),
	)

	log.Info("Initializing get all stores process")

	var stores []*domain.Store
	if err := s.db.WithContext(ctx).Where("userId = ?", userID.String()).Find(&stores).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("store not found")
			return nil, nil
		}

		log.Error("Failed to stores", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("stores found successfully")

	return stores, nil
}

func (s *storeRepository) GetByID(ctx context.Context, ID uuid.UUID) (*domain.Store, error) {
	log := slog.With(
		slog.String("func", "GetByID"),
		slog.String("repository", "store"),
	)

	log.Info("Initializing get store by id process")

	var store *domain.Store
	if err := s.db.WithContext(ctx).Where("id = ?", ID.String()).First(&store).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("store not found")
			return nil, nil
		}

		log.Error("Failed to stores", slog.String("error", err.Error()))
		return nil, err
	}
	log.Info("get store by id excuted sucessfully")
	return store, nil
}

func (s *storeRepository) UpdateName(ctx context.Context, name string, storeID uuid.UUID) error {
	log := slog.With(
		slog.String("func", "UpdateName"),
		slog.String("repository", "store"),
	)

	log.Info("Initializing update store name process")

	if err := s.db.WithContext(ctx).Model(domain.Store{}).Where("id = ?", storeID.String()).Update("name", name).Error; err != nil {
		log.Error("Failed to update store name", slog.String("error", err.Error()))
		return err
	}

	log.Info("store name updated successfully")
	return nil
}

func (s *storeRepository) Delete(ctx context.Context, storeID uuid.UUID) error {
	log := slog.With(
		slog.String("repository", "store"),
		slog.String("func", "Delete"),
	)

	log.Info("Initializing delete store process")

	if err := s.db.WithContext(ctx).Where("id = ?", storeID.String()).Delete(&domain.Store{}).Error; err != nil {
		log.Error("Failed to delete store", slog.String("error", err.Error()))
		return err
	}

	log.Info("store deleted successfully")
	return nil
}
