package domain

import (
	"context"
	"mime/multipart"
	"strings"
	"time"

	"github.com/GSVillas/e-commercer-api/util"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Billboard struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey;column:id"`
	StoreID   uuid.UUID      `gorm:"type:char(36);column:storeId;not null"`
	Store     Store          `gorm:"foreignKey:StoreID"`
	Label     string         `gorm:"size:100;not null;column:label"`
	ImageURL  string         `gorm:"size:400;not null;column:imageUrl"`
	CreatedAt time.Time      `gorm:"column:createdAt"`
	UpdatedAt time.Time      `gorm:"column:updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deletedAt"`
}

type BillboardPayload struct {
	Label string                `json:"label" validate:"required,min=1,max=100"`
	Image *multipart.FileHeader `json:"image" validate:"required"`
}

type BillboardRespose struct {
	ID        string    `json:"id"`
	Label     string    `json:"label"`
	StoreID   string    `json:"storeId"`
	ImageURL  string    `json:"imageUrl"`
	CreatedAt time.Time `json:"createdAt"`
}

type BillboardHandler interface {
	Create(ctx echo.Context) error
}

type BillboardService interface {
	Create(ctx context.Context, storeID uuid.UUID, billboardPayload BillboardPayload) (*BillboardRespose, error)
}

type BillboardRepository interface {
	Create(ctx context.Context, billboard Billboard) error
}

func (b *BillboardPayload) trim() {
	b.Label = strings.TrimSpace(b.Label)
}

func (b *BillboardPayload) Validate() error {
	b.trim()
	validate := validator.New()
	if err := validate.Struct(b); err != nil {
		return err
	}

	if err := util.ValidateFile(b.Image); err != nil {
		return err
	}

	return nil
}

func (b *BillboardPayload) ToBillboard(imageURL string, StoreID uuid.UUID) *Billboard {
	return &Billboard{
		ID:        uuid.New(),
		Label:     b.Label,
		StoreID:   StoreID,
		ImageURL:  imageURL,
		CreatedAt: time.Now().UTC(),
	}
}

func (b *Billboard) ToResponse() *BillboardRespose {
	return &BillboardRespose{
		ID:        b.ID.String(),
		Label:     b.Label,
		StoreID:   b.StoreID.String(),
		ImageURL:  b.ImageURL,
		CreatedAt: b.CreatedAt,
	}
}

func (Billboard) TableName() string {
	return "Billboard"
}
