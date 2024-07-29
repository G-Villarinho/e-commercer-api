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

type userHandler struct {
	i           *do.Injector
	userService domain.UserService
}

func NewUserHandler(i *do.Injector) (domain.UserHandler, error) {
	userService, err := do.Invoke[domain.UserService](i)
	if err != nil {
		return nil, err
	}

	return &userHandler{
		i:           i,
		userService: userService,
	}, nil
}

func (u *userHandler) Create(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "user"),
		slog.String("func", "Create"),
	)

	log.Info("Initializing user creation process")

	var userPayload domain.UserPayLoad
	if err := ctx.Bind(&userPayload); err != nil {
		log.Warn("Failed to bind payload", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusUnprocessableEntity, &problem.ProblemDetail{
			Status: http.StatusUnprocessableEntity,
			Title:  "Invalid Request",
			Detail: "Oops! Something went wrong while processing your request. Please try again later.",
		})
	}

	if err := userPayload.Validate(); err != nil {
		log.Warn("Invalid payload", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusBadRequest, &problem.ProblemDetail{
			Status: http.StatusBadRequest,
			Title:  "Invalid Request",
			Detail: "The data provided is incorrect or incomplete. Please verify and try again.",
		})
	}

	if err := u.userService.Create(ctx.Request().Context(), userPayload); err != nil {

		if errors.Is(err, domain.ErrUserAlreadyExists) {
			log.Warn("User already exists")
			return ctx.JSON(http.StatusConflict, &problem.ProblemDetail{
				Status: http.StatusConflict,
				Title:  "Conflict",
				Detail: "The user already exists. Please try again with a different email.",
			})
		}

		if errors.Is(err, domain.ErrHashingPassword) {
			log.Error("Failed to hash password", slog.String("error", err.Error()))
			return ctx.JSON(http.StatusInternalServerError, &problem.ProblemDetail{
				Status: http.StatusInternalServerError,
				Title:  "Internal Server Error",
				Detail: "Oops! Something went wrong while processing your request. Please try again later.",
			})
		}

		log.Error("Failed to create user", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusInternalServerError, &problem.ProblemDetail{
			Status: http.StatusInternalServerError,
			Title:  "Internal Server Error",
			Detail: "Oops! Something went wrong while processing your request. Please try again later.",
		})
	}

	log.Info("User created successfully")
	return ctx.NoContent(http.StatusCreated)
}

func (u *userHandler) SignIn(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "user"),
		slog.String("func", "SignIn"),
	)

	log.Info("Initializing user sign-in process")

	var signInPayload domain.SignInPayLoad
	if err := ctx.Bind(&signInPayload); err != nil {
		log.Warn("Failed to bind payload", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusUnprocessableEntity, &problem.ProblemDetail{
			Status: http.StatusUnprocessableEntity,
			Title:  "Invalid Request",
			Detail: "Oops! Something went wrong while processing your request. Please try again later.",
		})
	}

	if err := signInPayload.Validate(); err != nil {
		log.Warn("Invalid payload", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusBadRequest, &problem.ProblemDetail{
			Status: http.StatusBadRequest,
			Title:  "Invalid Request",
			Detail: "The data provided is incorrect or incomplete. Please verify and try again.",
		})
	}

	signInResponse, err := u.userService.SignIn(ctx.Request().Context(), signInPayload)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) || errors.Is(err, domain.ErrInvalidPassword) {
			log.Warn("Invalid credentials")
			return ctx.JSON(http.StatusUnauthorized, &problem.ProblemDetail{
				Status: http.StatusUnauthorized,
				Title:  "Unauthorized",
				Detail: "Invalid credentials. Please verify your email and password and try again.",
			})
		}

		if errors.Is(err, domain.ErrEmailNotConfirmed) {
			log.Warn("Email not confirmed")
			return ctx.JSON(http.StatusConflict, signInResponse)
		}

		log.Error("Failed to sign in user", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusInternalServerError, &problem.ProblemDetail{
			Status: http.StatusInternalServerError,
			Title:  "Internal Server Error",
			Detail: "Oops! Something went wrong while processing your request. Please try again later.",
		})
	}

	log.Info("User signed in successfully")
	return ctx.JSON(http.StatusOK, signInResponse)
}

func (u *userHandler) UpdateName(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "user"),
		slog.String("func", "UpdateName"),
	)

	log.Info("Initializing user name update process")

	var userUpdateNamePayload domain.UpdateNamePayload
	if err := ctx.Bind(&userUpdateNamePayload); err != nil {
		log.Warn("Failed to bind payload", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusUnprocessableEntity, &problem.ProblemDetail{
			Status: http.StatusUnprocessableEntity,
			Title:  "Invalid Request",
			Detail: "Oops! Something went wrong while processing your request. Please try again later.",
		})
	}

	if err := userUpdateNamePayload.Validate(); err != nil {
		log.Warn("Invalid payload", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusBadRequest, &problem.ProblemDetail{
			Status: http.StatusBadRequest,
			Title:  "Invalid Request",
			Detail: "The data provided is incorrect or incomplete. Please verify and try again.",
		})
	}

	if err := u.userService.UpdateName(ctx.Request().Context(), userUpdateNamePayload.Name); err != nil {
		if errors.Is(err, domain.ErrNameIsSame) {
			log.Warn("Name is same as before")
			return ctx.JSON(http.StatusConflict, &problem.ProblemDetail{
				Status: http.StatusConflict,
				Title:  "Conflict",
				Detail: "The name provided is same as before. Please try again with a different name.",
			})
		}

		if errors.Is(err, domain.ErrEmailNotConfirmed) {
			log.Warn("Email not confirmed")
			return ctx.JSON(http.StatusForbidden, &problem.ProblemDetail{
				Status: http.StatusForbidden,
				Title:  "Email Not Confirmed",
				Detail: "Your email address is not confirmed. Please check your inbox for the confirmation email and follow the instructions to confirm your email address.",
			})
		}

		log.Error("Failed to update user name", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusInternalServerError, &problem.ProblemDetail{
			Status: http.StatusInternalServerError,
			Title:  "Internal Server Error",
			Detail: "Oops! Something went wrong while processing your request. Please try again later.",
		})
	}

	log.Info("User name updated successfully")
	return ctx.NoContent(http.StatusOK)
}

func (u *userHandler) UpdatePassword(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "user"),
		slog.String("func", "UpdatePassword"),
	)

	log.Info("Initializing user password update process")

	var updatePasswordPayload domain.UpdatePasswordPayload
	if err := ctx.Bind(&updatePasswordPayload); err != nil {
		log.Warn("Failed to bind payload", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusUnprocessableEntity, &problem.ProblemDetail{
			Status: http.StatusUnprocessableEntity,
			Title:  "Invalid Request",
			Detail: "Oops! Something went wrong while processing your request. Please try again later.",
		})
	}

	if err := updatePasswordPayload.Validate(); err != nil {
		log.Warn("Invalid payload", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusBadRequest, &problem.ProblemDetail{
			Status: http.StatusBadRequest,
			Title:  "Invalid Request",
			Detail: "The data provided is incorrect or incomplete. Please verify and try again.",
		})
	}

	if err := u.userService.UpdatePassword(ctx.Request().Context(), updatePasswordPayload); err != nil {

		if errors.Is(err, domain.ErrEmailNotConfirmed) {
			log.Warn("Email not confirmed")
			return ctx.JSON(http.StatusForbidden, &problem.ProblemDetail{
				Status: http.StatusForbidden,
				Title:  "Email Not Confirmed",
				Detail: "Your email address is not confirmed. Please check your inbox for the confirmation email and follow the instructions to confirm your email address.",
			})
		}

		if errors.Is(err, domain.ErrInvalidOldPassword) {
			log.Warn("Invalid password")
			return ctx.JSON(http.StatusUnauthorized, &problem.ProblemDetail{
				Status: http.StatusUnauthorized,
				Title:  "Unauthorized",
				Detail: "Invalid old password. Please verify your password and try again.",
			})
		}

		if errors.Is(err, domain.ErrPasswordIsSame) {
			log.Warn("New password is same as old password")
			return ctx.JSON(http.StatusConflict, &problem.ProblemDetail{
				Status: http.StatusConflict,
				Title:  "Conflict",
				Detail: "The new password is same as the old password. Please try again with a different password.",
			})
		}

		log.Error("Failed to update user password", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusInternalServerError, &problem.ProblemDetail{
			Status: http.StatusInternalServerError,
			Title:  "Internal Server Error",
			Detail: "Oops! Something went wrong while processing your request. Please try again later.",
		})
	}

	log.Info("User password updated successfully")

	return ctx.NoContent(http.StatusOK)
}

func (u *userHandler) GetUserInfo(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "user"),
		slog.String("func", "GetUserInfos"),
	)

	log.Info("Initializing Get User Infos process")

	userResponse, err := u.userService.GetUserInfo(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &problem.ProblemDetail{
			Status: http.StatusInternalServerError,
			Title:  "Internal Server Error",
			Detail: "Oops! Something went wrong while processing your request. Please try again later.",
		})
	}

	log.Info("get user info executed successfully")
	return ctx.JSON(http.StatusOK, userResponse)
}

func (u *userHandler) ResendCode(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "user"),
		slog.String("func", "ResendCode"),
	)

	log.Info("Initializing resend code process")

	var resendCodePayload domain.ResendCodePayload

	if err := ctx.Bind(&resendCodePayload); err != nil {
		log.Warn("Failed to bind payload", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusUnprocessableEntity, &problem.ProblemDetail{
			Status: http.StatusUnprocessableEntity,
			Title:  "Invalid Request",
			Detail: "Oops! Something went wrong while processing your request. Please try again later.",
		})
	}

	if err := resendCodePayload.Validate(); err != nil {
		log.Warn("Invalid payload", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusBadRequest, &problem.ProblemDetail{
			Status: http.StatusBadRequest,
			Title:  "Invalid Request",
			Detail: "The data provided is incorrect or incomplete. Please verify and try again.",
		})
	}

	if err := u.userService.ResendCode(ctx.Request().Context(), resendCodePayload); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &problem.ProblemDetail{
			Status: http.StatusInternalServerError,
			Title:  "Internal Server Error",
			Detail: "Oops! Something went wrong while processing your request. Please try again later.",
		})
	}

	log.Info("code resend/send executed successfully")
	return ctx.NoContent(http.StatusOK)
}

func (u *userHandler) ConfirmEmail(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "user"),
		slog.String("func", "ConfirmEmail"),
	)

	log.Info("Initializing confirm email process")

	var confirmEmailPayload domain.ConfirmEmailPayload

	if err := ctx.Bind(&confirmEmailPayload); err != nil {
		log.Warn("Failed to bind payload", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusUnprocessableEntity, &problem.ProblemDetail{
			Status: http.StatusUnprocessableEntity,
			Title:  "Invalid Request",
			Detail: "Oops! Something went wrong while processing your request. Please try again later.",
		})
	}

	if err := confirmEmailPayload.Validate(); err != nil {
		log.Warn("Invalid payload", slog.String("error", err.Error()))
		return ctx.JSON(http.StatusBadRequest, &problem.ProblemDetail{
			Status: http.StatusBadRequest,
			Title:  "Invalid Request",
			Detail: "The data provided is incorrect or incomplete. Please verify and try again.",
		})
	}

	err := u.userService.ConfirmEmail(ctx.Request().Context(), confirmEmailPayload)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFoundInContext):
			log.Warn("User not found in context", slog.String("error", err.Error()))
			return ctx.JSON(http.StatusForbidden, &problem.ProblemDetail{
				Status: http.StatusForbidden,
				Title:  "Unauthorized",
				Detail: "You need to be logged in to access this resource.",
			})
		case errors.Is(err, domain.ErrOTPExpires):
			log.Warn("OTP has expired", slog.String("error", err.Error()))
			return ctx.JSON(http.StatusUnauthorized, &problem.ProblemDetail{
				Status: http.StatusUnauthorized,
				Title:  "OTP Expired",
				Detail: "The OTP code has expired. Please request a new code.",
			})
		case errors.Is(err, domain.ErrOTPInvalid):
			log.Warn("Invalid OTP", slog.String("error", err.Error()))
			return ctx.JSON(http.StatusUnauthorized, &problem.ProblemDetail{
				Status: http.StatusUnauthorized,
				Title:  "Invalid OTP",
				Detail: "The OTP code provided is invalid. Please check the code and try again.",
			})
		case errors.Is(err, domain.ErrSessionNotFound):
			log.Warn("Session not found", slog.String("error", err.Error()))
			return ctx.JSON(http.StatusUnauthorized, &problem.ProblemDetail{
				Status: http.StatusUnauthorized,
				Title:  "Session Not Found",
				Detail: "The session was not found. Please log in and try again.",
			})
		case errors.Is(err, domain.ErrTokenInvalid):
			log.Warn("Invalid token", slog.String("error", err.Error()))
			return ctx.JSON(http.StatusUnauthorized, &problem.ProblemDetail{
				Status: http.StatusUnauthorized,
				Title:  "Invalid Token",
				Detail: "The token provided is invalid. Please log in and try again.",
			})
		default:
			log.Error("Failed to confirm email", slog.String("error", err.Error()))
			return ctx.JSON(http.StatusInternalServerError, &problem.ProblemDetail{
				Status: http.StatusInternalServerError,
				Title:  "Internal Server Error",
				Detail: "Oops! Something went wrong while processing your request. Please try again later.",
			})
		}
	}

	log.Info("Email confirmed successfully")
	return ctx.NoContent(http.StatusOK)
}
