package http

import (
	"qttf/internal/player"

	"github.com/labstack/echo/v4"
)

func MapPlayerRoutes(playerGroup *echo.Group, h player.Handlers) {
	playerGroup.GET("/", h.GetPlayers())
	playerGroup.GET("/:player_id", h.GetPlayerById())
}
