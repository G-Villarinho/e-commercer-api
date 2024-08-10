package database

import (
	"context"

	"github.com/GSVillas/e-commercer-api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConnection(ctx context.Context) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.Env.ConnectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	return db, nil
}
