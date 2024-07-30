package domain

import (
	"time"

	"github.com/google/uuid"
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
	Label    string `json:"label" validate:"required,min=1,max=100"`
	ImageURL string `json:"imageUrl" validate:"required,min=1,max=400"`
}

func (Billboard) TableName() string {
	return "Billboard"
}
