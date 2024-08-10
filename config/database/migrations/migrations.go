package main

import (
	"context"
	"log"
	"time"

	"github.com/GSVillas/e-commercer-api/config"
	"github.com/GSVillas/e-commercer-api/config/database"
	"github.com/GSVillas/e-commercer-api/domain"
)

func main() {
	config.LoadEnvironments()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.NewMysqlConnection(ctx)
	if err != nil {
		log.Fatal("Fail to connect to mysql: ", err)
	}

	if err := db.AutoMigrate(&domain.User{}, &domain.Store{}, &domain.Billboard{}); err != nil {
		log.Fatal("Fail to migrate: ", err)
	}

	log.Println("Migration executed successfully")
}
