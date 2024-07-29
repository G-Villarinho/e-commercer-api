package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/OVillas/e-commercer-api/domain"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/problem-details"
	"github.com/samber/do"
)

type storeHandler struct {
	i            *do.Injector
	storeService domain.StoreService
}

func NewStoreHandler(i *do.Injector) (domain.StoreHandler, error) {
	storeService, err := do.Invoke[domain.StoreService](i)
	if err != nil {
		return nil, err
	}

	return &storeHandler{
		i:            i,
		storeService: storeService,
	}, nil
}

func (s *storeHandler) Create(ctx echo.Context) error {
	log := slog.With(
		slog.String("func", "Create"),
		slog.String("handler", "store"),
	)

	log.Info("Initializing store create process")

	var storePayload domain.StorePayload
	if err := ctx.Bind(&storePayload); err != nil {
		log.Warn("Failed to bind payload", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusUnprocessableEntity, &problem.ProblemDetail{
			Status: http.StatusUnprocessableEntity,
			Title:  "Invalid Request",
			Detail: "Oops! Something went wrong while processing your request. Please try again later.",
		})
	}

	if err := storePayload.Validate(); err != nil {
		log.Warn("Invalid payload", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusBadRequest, &problem.ProblemDetail{
			Status: http.StatusBadRequest,
			Title:  "Invalid Request",
			Detail: "The data provided is incorrect or incomplete. Please verify and try again.",
		})
	}

	storeResponse, err := s.storeService.Create(ctx.Request().Context(), storePayload)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &problem.ProblemDetail{
			Status: http.StatusInternalServerError,
			Title:  "Internal Server Error",
			Detail: "Oops! Something went wrong while processing your request. Please try again later.",
		})
	}

	log.Info("Store created successfully")

	return ctx.JSON(http.StatusCreated, storeResponse)
}

func (s *storeHandler) GetAll(ctx echo.Context) error {
	log := slog.With(
		slog.String("func", "GetAll"),
		slog.String("handler", "store"),
	)

	log.Info("Initializing get all stores process")

	storeResponse, err := s.storeService.GetAll(ctx.Request().Context())
	if err != nil {
		log.Error("Failed to get all stores", slog.String("error", err.Error()))

		switch {
		case errors.Is(err, domain.ErrUserNotFoundInContext):
			return ctx.JSON(http.StatusForbidden, &problem.ProblemDetail{
				Status: http.StatusForbidden,
				Title:  "Forbidden",
				Detail: "User not found in context. Please log in again.",
			})
		case errors.Is(err, domain.ErrStoresNotFound):
			return ctx.NoContent(http.StatusNoContent)
		default:
			return ctx.JSON(http.StatusInternalServerError, &problem.ProblemDetail{
				Status: http.StatusInternalServerError,
				Title:  "Internal Server Error",
				Detail: "Oops! Something went wrong while processing your request. Please try again later.",
			})
		}
	}

	log.Info("Successfully retrieved all stores")
	return ctx.JSON(http.StatusOK, storeResponse)
}
