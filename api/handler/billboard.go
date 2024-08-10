package handler

import (
	"log/slog"
	"net/http"

	"github.com/GSVillas/e-commercer-api/domain"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/problem-details"
	"github.com/samber/do"
)

type billboardHandler struct {
	i                *do.Injector
	billboardService domain.BillboardService
	userService      domain.UserService
}

func NewBillboardHandler(i *do.Injector) (domain.BillboardHandler, error) {
	billboardService, err := do.Invoke[domain.BillboardService](i)
	if err != nil {
		return nil, err
	}

	userService, err := do.Invoke[domain.UserService](i)
	if err != nil {
		return nil, err
	}

	return &billboardHandler{
		i:                i,
		billboardService: billboardService,
		userService:      userService,
	}, nil
}

func (b *billboardHandler) Create(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "billboard"),
		slog.String("func", "Create"),
	)

	log.Info("Initializing billboard create process")

	param := ctx.Param("storeId")

	storeID, err := uuid.Parse(param)
	if err != nil {
		log.Warn("Invalid params", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusBadRequest, &problem.ProblemDetail{
			Status: http.StatusBadRequest,
			Title:  "Invalid Request",
			Detail: "The data provided is incorrect or incomplete. Please verify and try again.",
		})
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		log.Warn("Image file is invalid", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusBadRequest, &problem.ProblemDetail{
			Status: http.StatusBadRequest,
			Title:  "Invalid Request",
			Detail: "The 'image' field is invalid.",
		})
	}

	billboardPayload := &domain.BillboardPayload{
		Label: ctx.FormValue("label"),
		Image: file,
	}

	if err = billboardPayload.Validate(); err != nil {
		log.Warn("Invalid payload", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusBadRequest, &problem.ProblemDetail{
			Status: http.StatusBadRequest,
			Title:  "Invalid Request",
			Detail: "The data provided is incorrect or incomplete. Please verify and try again.",
		})
	}

	billboardResponse, err := b.billboardService.Create(ctx.Request().Context(), storeID, *billboardPayload)
	if err != nil {
		switch err {
		case domain.ErrUserNotFoundInContext:
			log.Warn("User not found in context", slog.String("error", err.Error()))
			return ctx.JSON(http.StatusUnauthorized, &problem.ProblemDetail{
				Status: http.StatusUnauthorized,
				Title:  "Unauthorized",
				Detail: "User not authorized to perform this action.",
			})
		case domain.ErrStoreNotFound:
			log.Warn("Store not found", slog.String("error", err.Error()))
			return ctx.JSON(http.StatusNotFound, &problem.ProblemDetail{
				Status: http.StatusNotFound,
				Title:  "Store Not Found",
				Detail: "The specified store was not found.",
			})
		case domain.ErrUnauthorizedAction:
			log.Warn("Unauthorized action attempted", slog.String("error", err.Error()))
			return ctx.JSON(http.StatusForbidden, &problem.ProblemDetail{
				Status: http.StatusForbidden,
				Title:  "Forbidden",
				Detail: "You are not allowed to perform this action.",
			})
		default:
			log.Error("Unexpected error", slog.String("error", err.Error()))
			return ctx.JSON(http.StatusInternalServerError, &problem.ProblemDetail{
				Status: http.StatusInternalServerError,
				Title:  "Internal Server Error",
				Detail: "Oops! Something went wrong while processing your request. Please try again later.",
			})
		}
	}

	return ctx.JSON(http.StatusOK, billboardResponse)
}
