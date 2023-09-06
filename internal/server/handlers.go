package server

import (
	"net/http"
	cityHttp "qttf/internal/city/delivery/http"
	cityRepository "qttf/internal/city/repository"
	cityUsecase "qttf/internal/city/usecase"
	playerHttp "qttf/internal/player/delivery/http"
	playerRepository "qttf/internal/player/repository"
	playerUsecase "qttf/internal/player/usecase"
	ratingHttp "qttf/internal/rating/delivery/http"
	ratingRepository "qttf/internal/rating/repository"
	ratingUsecase "qttf/internal/rating/usecase"

	"github.com/labstack/echo/v4"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	// Init repositories
	cityRepo := cityRepository.NewCityRepository(s.db)
	playerRepo := playerRepository.NewPlayerRepository(s.db)
	ratingRepo := ratingRepository.NewRatingRepository(s.db)

	// Init usecases
	cityUC := cityUsecase.NewCityUseCase(cityRepo)
	playerUC := playerUsecase.NewPlayerUseCase(playerRepo)
	ratingUC := ratingUsecase.NewRatingUseCase(ratingRepo)

	// Init handlers
	cityHandlers := cityHttp.NewCityHandlers(cityUC)
	playerHandlers := playerHttp.NewPlayerHandlers(playerUC)
	ratingHandlers := ratingHttp.NewRatingHandlers(ratingUC)

	v1 := e.Group("/api/v1")

	health := v1.Group("/health")
	cityGroup := v1.Group("/cities")
	playerGroup := v1.Group("/players")
	ratingGroup := v1.Group("/rating")

	cityHttp.MapCityRoutes(cityGroup, cityHandlers)
	playerHttp.MapPlayerRoutes(playerGroup, playerHandlers)
	ratingHttp.MapRatingRoutes(ratingGroup, ratingHandlers)

	health.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return nil
}
