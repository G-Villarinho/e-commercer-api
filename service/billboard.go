package service

import (
	"context"
	"log/slog"

	"github.com/OVillas/e-commercer-api/client"
	"github.com/OVillas/e-commercer-api/domain"
	"github.com/OVillas/e-commercer-api/middleware"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type billboardService struct {
	i                   *do.Injector
	billboardRepository domain.BillboardRepository
	storeRepository     domain.StoreRepository
	cloudFlareService   client.CloudFlareService
}

func NewBillboardService(i *do.Injector) (domain.BillboardService, error) {
	billboardRepository, err := do.Invoke[domain.BillboardRepository](i)
	if err != nil {
		return nil, err
	}

	storeRepository, err := do.Invoke[domain.StoreRepository](i)
	if err != nil {
		return nil, err
	}

	cloudFlareService, err := do.Invoke[client.CloudFlareService](i)
	if err != nil {
		return nil, err
	}

	return &billboardService{
		i:                   i,
		billboardRepository: billboardRepository,
		storeRepository:     storeRepository,
		cloudFlareService:   cloudFlareService,
	}, nil
}

func (b *billboardService) Create(ctx context.Context, storeID uuid.UUID, billboardPayload domain.BillboardPayload) (*domain.BillboardRespose, error) {
	log := slog.With(
		slog.String("service", "billboard"),
		slog.String("func", "Create"),
	)

	log.Info("Initializing create billboard process")

	session, ok := ctx.Value(middleware.UserKey).(*domain.Session)
	if !ok || session == nil {
		return nil, domain.ErrUserNotFoundInContext
	}

	store, err := b.storeRepository.GetByID(ctx, storeID)
	if err != nil {
		log.Error("Failed to get store by id", slog.String("error", err.Error()))
		return nil, err
	}

	if store == nil {
		log.Warn("store not found with this id")
		return nil, domain.ErrStoreNotFound
	}

	if store.UserID != session.UserID {
		log.Error("Unauthorized attempt to delete store", slog.String("storeID", store.ID.String()), slog.String("userID", session.UserID.String()))
		return nil, domain.ErrUnauthorizedAction
	}

	imageURL, err := b.cloudFlareService.UploadImage(billboardPayload.Image)
	if err != nil {
		log.Error("Error to upload image in cloud", slog.String("error", err.Error()))
		return nil, err
	}

	billboard := billboardPayload.ToBillboard(imageURL, storeID)

	if err := b.billboardRepository.Create(ctx, *billboard); err != nil {
		log.Error("Error to create a billboard", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("Create billboard process executed succefully")
	return billboard.ToResponse(), nil
}
