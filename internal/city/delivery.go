package city

import "github.com/labstack/echo/v4"

type Handlers interface {
	GetCities() echo.HandlerFunc
}
