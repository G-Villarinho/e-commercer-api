package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/OVillas/e-commercer-api/api/handler"
	"github.com/OVillas/e-commercer-api/client"
	"github.com/OVillas/e-commercer-api/config"
	"github.com/OVillas/e-commercer-api/config/database"
	"github.com/OVillas/e-commercer-api/repository"
	"github.com/OVillas/e-commercer-api/service"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/resend/resend-go/v2"
	"github.com/samber/do"
	"gorm.io/gorm"
)

func main() {
	config.ConfigureLogger()
	config.LoadEnvironments()

	e := echo.New()
	i := do.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{config.Env.URLFront},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.NewMysqlConnection(ctx)
	if err != nil {
		log.Fatal("Fail to connect to mysql: ", err)
	}

	redisClient, err := database.NewRedisConnection(ctx)
	if err != nil {
		log.Fatal("Fail to connect to redis: ", err)
	}

	resendClient := resend.NewClient(config.Env.ResendKey)

	do.Provide(i, func(i *do.Injector) (*gorm.DB, error) {
		return db, nil
	})

	do.Provide(i, func(i *do.Injector) (*redis.Client, error) {
		return redisClient, nil
	})

	do.Provide(i, func(i *do.Injector) (*resend.Client, error) {
		return resendClient, nil
	})

	do.Provide(i, handler.NewHealthCheckHandler)
	do.Provide(i, handler.NewUserHandler)
	do.Provide(i, service.NewUserService)
	do.Provide(i, repository.NewUserRepository)
	do.Provide(i, service.NewSessionService)
	do.Provide(i, repository.NewSessionRepository)
	do.Provide(i, handler.NewStoreHandler)
	do.Provide(i, service.NewStoreService)
	do.Provide(i, repository.NewStoreRepository)
	do.Provide(i, service.NewEmailService)
	do.Provide(i, client.NewCloudFlareService)
	do.Provide(i, handler.NewBillboardHandler)
	do.Provide(i, service.NewBillboardService)
	do.Provide(i, repository.NewBillboardRepository)

	handler.SetupRoutes(e, i)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Env.APIPort)))
}
