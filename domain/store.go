package domain

import (
	"context"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Store struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;column:id"`
	Name      string         `gorm:"size:100;not null;column:name"`
	UserID    uuid.UUID      `gorm:"column:userId;not null"`
	User      User           `gorm:"foreignKey:UserID"`
	CreatedAt time.Time      `gorm:"column:createdAt"`
	UpdatedAt time.Time      `gorm:"column:updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deletedAt"`
}

type StorePayload struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

type StoreResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type StoreHandler interface {
	Create(ctx echo.Context) error
}

type StoreService interface {
	Create(ctx context.Context, storePayload StorePayload) error
}

type StoreRepository interface {
	Create(ctx context.Context, store Store) error
}

func (s *StorePayload) trim() {
	s.Name = strings.TrimSpace(s.Name)
}

func (s *StorePayload) Validate() error {
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
		UpdatedAt: time.Now().UTC(),
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
