package main

import (
	"context"
	"log"
	"time"

	"github.com/OVillas/e-commercer-api/config"
	"github.com/OVillas/e-commercer-api/config/database"
	"github.com/OVillas/e-commercer-api/domain"
)

func main() {
	config.LoadEnvironments()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.NewPostgresConnection(ctx)
	if err != nil {
		log.Fatal("Fail to connect to postgres: ", err)
	}

	if err := db.AutoMigrate(&domain.User{}, &domain.Store{}); err != nil {
		log.Fatal("Fail to migrate: ", err)
	}

	log.Println("Migration executed successfully")
}
