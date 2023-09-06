package rating

import "github.com/labstack/echo/v4"

type Handlers interface {
	GetRating() echo.HandlerFunc
}
