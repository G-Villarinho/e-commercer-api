package domain

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var (
	ErrStoresNotFound     = errors.New("stores not found for this userID")
	ErrStoreNotFound      = errors.New("store not found for this userID")
	ErrUnauthorizedAction = errors.New("unauthorized action")
)

type Store struct {
	ID         uuid.UUID      `gorm:"type:char(36);primaryKey;column:id"`
	Name       string         `gorm:"size:100;not null;column:name"`
	UserID     uuid.UUID      `gorm:"type:char(36);column:userId;not null"`
	User       User           `gorm:"foreignKey:UserID"`
	Billboards []Billboard    `gorm:"foreignKey:StoreID"`
	CreatedAt  time.Time      `gorm:"column:createdAt"`
	UpdatedAt  time.Time      `gorm:"column:updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index;column:deletedAt"`
}

type StorePayload struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}
type StoreNameUpdatePayload struct {
	ID   uuid.UUID `json:"id" validate:"required"`
	Name string    `json:"name" validate:"required,min=1,max=100"`
}

type StoreResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type StoreHandler interface {
	Create(ctx echo.Context) error
	GetAll(ctx echo.Context) error
	UpdateName(ctx echo.Context) error
}

type StoreService interface {
	Create(ctx context.Context, storePayload StorePayload) (*StoreResponse, error)
	GetAll(ctx context.Context) ([]*StoreResponse, error)
	UpdateName(ctx context.Context, updateStoreNamePayload StoreNameUpdatePayload) error
}

type StoreRepository interface {
	Create(ctx context.Context, store Store) error
	GetAll(ctx context.Context, userID uuid.UUID) ([]*Store, error)
	GetByID(ctx context.Context, ID uuid.UUID) (*Store, error)
	UpdateName(ctx context.Context, name string, ID uuid.UUID) error
}

func (s *StorePayload) trim() {
	s.Name = strings.TrimSpace(s.Name)
}

func (s *StoreNameUpdatePayload) trim() {
	s.Name = strings.TrimSpace(s.Name)
}

func (s *StorePayload) Validate() error {
	s.trim()
	validator := validator.New()
	return validator.Struct(s)
}

func (s *StoreNameUpdatePayload) Validate() error {
	s.trim()
	validator := validator.New()
	return validator.Struct(s)
}

func (s *StorePayload) ToStore(userID uuid.UUID) *Store {
	return &Store{
		ID:        uuid.New(),
		Name:      s.Name,
		UserID:    userID,
		CreatedAt: time.Now().UTC(),
	}
}

func (s *Store) ToResponse() *StoreResponse {
	return &StoreResponse{
		ID:        s.ID.String(),
		Name:      s.Name,
		CreatedAt: s.CreatedAt.Format(time.RFC3339),
	}
}

func (Store) TableName() string {
	return "Store"
}
