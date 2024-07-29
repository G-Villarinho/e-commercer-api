package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/OVillas/e-commercer-api/config"
	"github.com/OVillas/e-commercer-api/domain"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/problem-details"
	"github.com/samber/do"
)

type contextKey string

const UserKey contextKey = "user"

func CheckLoggedIn(i *do.Injector) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			userSession := do.MustInvoke[domain.SessionService](i)

			authorizationHeader := ctx.Request().Header.Get("Authorization")
			if authorizationHeader == "" {
				return ctx.JSON(http.StatusUnauthorized, &problem.ProblemDetail{
					Status: http.StatusUnauthorized,
					Title:  "Access Denied",
					Detail: "You need to be logged in to access this resource.",
				})
			}

			tokenString, err := extractToken(ctx)
			if err != nil {
				return ctx.JSON(http.StatusUnauthorized, &problem.ProblemDetail{
					Status: http.StatusUnauthorized,
					Title:  "Invalid Session",
					Detail: "Your session is invalid or missing. Please log in again.",
				})
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
					return nil, domain.ErrorUnexpectedMethod
				}
				return config.Env.PublicKey, nil
			})

			if err != nil || !token.Valid {
				return ctx.JSON(http.StatusUnauthorized, &problem.ProblemDetail{
					Status: http.StatusForbidden,
					Title:  "Invalid Session",
					Detail: "Your session is invalid. Please log in again.",
				})
			}

			user, err := userSession.GetUser(ctx.Request().Context(), tokenString)
			if err != nil {
				if errors.Is(err, domain.ErrSessionNotFound) {
					return ctx.JSON(http.StatusForbidden, &problem.ProblemDetail{
						Status: http.StatusUnauthorized,
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

func extractToken(ctx echo.Context) (string, error) {
	token := ctx.Request().Header.Get("Authorization")

	length := len(strings.Split(token, " "))
	if length == 2 {
		return strings.Split(token, " ")[1], nil
	}

	return "", domain.ErrSessionNotFound
}
