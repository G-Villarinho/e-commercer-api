package handler

import (
	"log/slog"
	"net/http"

	"github.com/OVillas/e-commercer-api/domain"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/problem-details"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type healthCheckHandler struct {
	i           *do.Injector
	db          *gorm.DB
	redisClient *redis.Client
}

func NewHealthCheckHandler(i *do.Injector) (domain.HealthCheckHandler, error) {
	db := do.MustInvoke[*gorm.DB](i)
	redisClient := do.MustInvoke[*redis.Client](i)

	return &healthCheckHandler{
		i:           i,
		db:          db,
		redisClient: redisClient,
	}, nil
}

func (h *healthCheckHandler) HealthCheck(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "HealthCheck"),
		slog.String("func", "HealthCheck"),
	)

	log.Info("HealthCheck initiated")

	context := ctx.Request().Context()

	_, err := h.redisClient.Ping(context).Result()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &problem.ProblemDetail{
			Status: http.StatusInternalServerError,
			Title:  "Internal Server Error",
			Detail: "Oops! Something went wrong while processing your request. Please try again later.",
		})
	}

	sqlDB, err := h.db.DB()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &problem.ProblemDetail{
			Status: http.StatusInternalServerError,
			Title:  "Internal Server Error",
			Detail: "Oops! Something went wrong while processing your request. Please try again later.",
		})
	}

	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return ctx.JSON(http.StatusInternalServerError, &problem.ProblemDetail{
			Status: http.StatusInternalServerError,
			Title:  "Internal Server Error",
			Detail: "Oops! Something went wrong while processing your request. Please try again later.",
		})
	}

	return ctx.NoContent(http.StatusOK)
}
