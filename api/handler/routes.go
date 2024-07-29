package handler

import (
	"github.com/OVillas/e-commercer-api/domain"
	Middleware "github.com/OVillas/e-commercer-api/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samber/do"
)

func SetupRoutes(e *echo.Echo, i *do.Injector) {
	setupHealthCheckRoutes(e, i)
	setupUserRoutes(e, i)
	setupStoreRoutes(e, i)
}

func setupHealthCheckRoutes(e *echo.Echo, i *do.Injector) {
	healthCheckHandler := do.MustInvoke[domain.HealthCheckHandler](i)
	e.GET("/health", healthCheckHandler.HealthCheck)
}

func setupUserRoutes(e *echo.Echo, i *do.Injector) {
	userHandler := do.MustInvoke[domain.UserHandler](i)
	group := e.Group("/v1/users")
	group.POST("", userHandler.Create)
	group.POST("/signIn", userHandler.SignIn)
	group.PATCH("/name", userHandler.UpdateName, Middleware.CheckLoggedIn(i))
	group.PATCH("/password", userHandler.UpdatePassword, Middleware.CheckLoggedIn(i))
	group.GET("/me", userHandler.GetUserInfo, Middleware.CheckLoggedIn(i))
	group.PATCH("/email/confirm", userHandler.ConfirmEmail, Middleware.CheckLoggedIn(i))
	group.POST("/resend-code", userHandler.ResendCode, middleware.RateLimiterWithConfig(
		middleware.RateLimiterConfig{
			Skipper: middleware.DefaultSkipper,
			Store:   middleware.NewRateLimiterMemoryStore(20),
		}))
}

func setupStoreRoutes(e *echo.Echo, i *do.Injector) {
	storeHandler := do.MustInvoke[domain.StoreHandler](i)
	group := e.Group("/v1/stores", Middleware.CheckLoggedIn(i))
	group.POST("", storeHandler.Create)
	group.GET("", storeHandler.GetAll)
}
