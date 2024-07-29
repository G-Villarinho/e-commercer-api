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

func (s *storeService) Create(ctx context.Context, storePayload domain.StorePayload) (*domain.StoreResponse, error) {
	log := slog.With(
		slog.String("service", "store"),
		slog.String("func", "Create"),
	)

	log.Info("Initializing store creation process")

	session, ok := ctx.Value(middleware.UserKey).(*domain.Session)
	if !ok || session == nil {
		return nil, domain.ErrUserNotFoundInContext
	}

	store := storePayload.ToStore(session.UserID)

	if err := s.storeRepository.Create(ctx, *store); err != nil {
		log.Error("Failed to create store", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("store creation process executed sucessfully")
	return store.ToResponse(), nil
}

func (s *storeService) GetAll(ctx context.Context) ([]*domain.StoreResponse, error) {
	log := slog.With(
		slog.String("service", "store"),
		slog.String("func", "GetAll"),
	)

	log.Info("Initializing store retrieval process")

	session, ok := ctx.Value(middleware.UserKey).(*domain.Session)
	if !ok || session == nil {
		log.Error("User not found in context")
		return nil, domain.ErrUserNotFoundInContext
	}

	stores, err := s.storeRepository.GetAll(ctx, session.UserID)
	if err != nil {
		log.Error("Failed to retrieve stores", slog.String("error", err.Error()))
		return nil, err
	}

	if stores == nil {
		log.Warn("No stores found for the user", slog.String("userID", session.UserID.String()))
		return nil, domain.ErrStoresNotFound
	}

	log.Info("Successfully retrieved stores", slog.Int("storeCount", len(stores)))

	var storesResponse []*domain.StoreResponse
	for _, store := range stores {
		storesResponse = append(storesResponse, store.ToResponse())
	}

	log.Info("Store retrieval process executed successfully")
	return storesResponse, nil
}
