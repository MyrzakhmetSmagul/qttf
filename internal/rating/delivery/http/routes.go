package http

import (
	"qttf/internal/rating"

	"github.com/labstack/echo/v4"
)

func MapRatingRoutes(ratingGroup *echo.Group, h rating.Handlers) {
	ratingGroup.GET("/", h.GetRating())
}
