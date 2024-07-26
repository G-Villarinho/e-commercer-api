package service

import (
	"context"
	"log/slog"

	"github.com/OVillas/e-commercer-api/domain"
	"github.com/OVillas/e-commercer-api/middleware"
	"github.com/samber/do"
)

type storeService struct {
	i               *do.Injector
	storeRepository domain.StoreRepository
}

func NewStoreService(i *do.Injector) (domain.StoreService, error) {
	storeRepository, err := do.Invoke[domain.StoreRepository](i)
	if err != nil {
		return nil, err
	}

	return &storeService{
		i:               i,
		storeRepository: storeRepository,
	}, nil
}

func (s *storeService) Create(ctx context.Context, storePayload domain.StorePayload) error {
	log := slog.With(
		slog.String("service", "store"),
		slog.String("func", "Create"),
	)

	log.Info("Initializing store creation process")

	user, ok := ctx.Value(middleware.UserKey).(*domain.UserSession)
	if !ok || user == nil {
		return domain.ErrUserNotFoundInContext
	}

	store := storePayload.ToStore(user.UserID)

	if err := s.storeRepository.Create(ctx, *store); err != nil {
		log.Error("Failed to create store", slog.String("error", err.Error()))
		return err
	}

	return nil
}
