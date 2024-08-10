package repository

import (
	"context"
	"log/slog"

	"github.com/GSVillas/e-commercer-api/domain"
	"github.com/go-redis/redis/v8"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type billboarRepository struct {
	i           *do.Injector
	db          *gorm.DB
	redisClient *redis.Client
}

func NewBillboardRepository(i *do.Injector) (domain.BillboardRepository, error) {
	db, err := do.Invoke[*gorm.DB](i)
	if err != nil {
		return nil, err
	}

	redisClient, err := do.Invoke[*redis.Client](i)
	if err != nil {
		return nil, err
	}

	return &billboarRepository{
		i:           i,
		db:          db,
		redisClient: redisClient,
	}, nil
}

func (b *billboarRepository) Create(ctx context.Context, billboard domain.Billboard) error {
	log := slog.With(
		slog.String("repository", "billboard"),
		slog.String("func", "Create"),
	)

	log.Info("Initializing billboard creation process")

	if err := b.db.WithContext(ctx).Create(&billboard).Error; err != nil {
		log.Error("Failed to create billboard", slog.String("error", err.Error()))
		return err
	}

	log.Info("billboard created successfully")
	return nil
}
