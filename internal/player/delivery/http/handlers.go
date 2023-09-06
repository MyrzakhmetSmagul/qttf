package http

import (
	"log"
	"net/http"
	"qttf/internal/player"

	echo "github.com/labstack/echo/v4"
)

type playerHandlers struct {
	playerUC player.UseCase
}

func NewPlayerHandlers(playerUC player.UseCase) player.Handlers {
	return &playerHandlers{playerUC: playerUC}
}

// GetPlayerById implements player.Handlers.
func (p *playerHandlers) GetPlayerById() echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO something
		return c.JSON(200, "in process...")
	}
}

// GetPlayers implements player.Handlers.
func (p *playerHandlers) GetPlayers() echo.HandlerFunc {
	return func(c echo.Context) error {
		players, err := p.playerUC.GetPlayers()
		if err != nil {
			log.Println("error playerHandlers.GetPlayers: ", err)
			return c.JSON(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, players)
	}
}
