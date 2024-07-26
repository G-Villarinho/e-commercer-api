package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/OVillas/e-commercer-api/config"
	"github.com/OVillas/e-commercer-api/domain"
	"github.com/OVillas/e-commercer-api/util"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/problem-details"
	"github.com/samber/do"
)

func ResendCode(i *do.Injector) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			userSession := do.MustInvoke[domain.UserSessionService](i)
			authorizationHeader := ctx.Request().Header.Get("Authorization")

			if authorizationHeader == "" {
				return next(ctx)
			}

			tokenString, err := util.ExtractToken(ctx)
			if err != nil {
				return ctx.JSON(http.StatusUnauthorized, &problem.ProblemDetail{
					Status: http.StatusUnauthorized,
					Title:  "Invalid Session",
					Detail: "Your session is invalid or missing. Please log in again.",
				})
			}

			publicKey, err := util.LoadPrivateKey(config.Env.SecretKeyPath)
			if err != nil {
				return ctx.JSON(http.StatusInternalServerError, &problem.ProblemDetail{
					Status: http.StatusInternalServerError,
					Title:  "Internal Server Error",
					Detail: "Failed to verify your account.",
				})
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
					return nil, domain.ErrorUnexpectedMethod
				}
				return publicKey, nil
			})

			if err != nil || !token.Valid {
				return ctx.JSON(http.StatusUnauthorized, &problem.ProblemDetail{
					Status: http.StatusUnauthorized,
					Title:  "Invalid Session",
					Detail: "Your session is invalid. Please log in again.",
				})
			}

			user, err := userSession.GetUser(ctx.Request().Context(), tokenString)
			if err != nil {
				if errors.Is(err, domain.ErrSessionNotFound) {
					return ctx.JSON(http.StatusForbidden, &problem.ProblemDetail{
						Status: http.StatusForbidden,
						Title:  "Session Expired",
						Detail: "Your session has expired. Please log in again to continue.",
					})
				}

				return ctx.JSON(http.StatusUnauthorized, &problem.ProblemDetail{
					Status: http.StatusUnauthorized,
					Title:  "Unauthorized Access",
					Detail: "Your session is invalid. Please log in again.",
				})
			}

			ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), UserKey, user)))
			return next(ctx)
		}
	}
}
