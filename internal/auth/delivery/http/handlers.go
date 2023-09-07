package http

import (
	"net/http"
	"qttf/internal/auth"

	echo "github.com/labstack/echo/v4"
)

type authHandlers struct {
	authUC auth.UseCase
}

func NewAuthHandlers(authUC auth.UseCase) auth.Handlers {
	return &authHandlers{authUC: authUC}
}

// SaveGoogleToken implements auth.Handlers.
func (a *authHandlers) SaveGoogleToken() echo.HandlerFunc {
	return func(c echo.Context) error {
		err := a.authUC.SaveGoogleToken(c.FormValue("code"))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.NoContent(http.StatusOK)
	}
}

// GetGoogleToken implements auth.Handlers.
func (a *authHandlers) GetGoogleToken() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, a.authUC.GetGoogleToken())
	}
}
