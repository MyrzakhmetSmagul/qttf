package http

import (
	"qttf/internal/auth"

	"github.com/labstack/echo/v4"
)

func MapAuthHandlers(authGroup *echo.Group, h auth.Handlers) {
	authGroup.GET("/", h.GetGoogleToken())
	authGroup.POST("/", h.SaveGoogleToken())
}
