package domain

import "github.com/labstack/echo/v4"

type HealthCheckHandler interface {
	HealthCheck(ctx echo.Context) error
}
