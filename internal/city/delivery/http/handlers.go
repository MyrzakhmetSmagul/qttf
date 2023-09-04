package http

import (
	"qttf/internal/city"

	"github.com/labstack/echo/v4"
)

type cityHandlers struct {
	cityRepo city.Repository
}

func (c *cityHandlers) GetCities() echo.HandlerFunc {
	return func(c echo.Context) error {

		return nil
	}
}
