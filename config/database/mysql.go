package database

import (
	"context"

	"github.com/OVillas/e-commercer-api/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysqlConnection(ctx context.Context) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(config.Env.ConnectionString), &gorm.Config{})
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
