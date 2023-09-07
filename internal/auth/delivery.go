package auth

import "github.com/labstack/echo/v4"

type Handlers interface {
	GetGoogleToken() echo.HandlerFunc
	SaveGoogleToken() echo.HandlerFunc
}
