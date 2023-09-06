package http

import (
	"qttf/internal/city"

	"github.com/labstack/echo/v4"
)

func MapCityRoutes(cityGroup *echo.Group, h city.Handlers) {
	cityGroup.GET("/", h.GetCities())
}
