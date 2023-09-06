package http

import (
	"log"
	"net/http"
	"qttf/internal/rating"

	echo "github.com/labstack/echo/v4"
)

type ratingHandlers struct {
	ratingUC rating.UseCase
}

// GetRating implements rating.Handlers.
func (r *ratingHandlers) GetRating() echo.HandlerFunc {
	return func(c echo.Context) error {
		rating, err := r.ratingUC.GetRating()
		if err != nil {
			log.Printf("ratingHandlers was failed: %v\n", err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, rating)
	}
}

func NewRatingHandlers(ratingUC rating.UseCase) rating.Handlers {
	return &ratingHandlers{ratingUC: ratingUC}
}
