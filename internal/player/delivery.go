package player

import (
	"github.com/labstack/echo/v4"
)

type Handlers interface {
	GetPlayers() echo.HandlerFunc
	GetPlayerById() echo.HandlerFunc
}
