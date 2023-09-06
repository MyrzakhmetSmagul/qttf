package http

import (
	"net/http"
	"qttf/internal/city"

	"github.com/labstack/echo/v4"
)

type cityHandlers struct {
	cityUC city.UseCase
}

func NewCityHandlers(cityUC city.UseCase) city.Handlers {
	return &cityHandlers{cityUC: cityUC}
}

func (ch *cityHandlers) GetCities() echo.HandlerFunc {
	return func(c echo.Context) error {
		cities, err := ch.cityUC.GetCities()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, cities)
	}
}
