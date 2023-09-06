package rating

import "qttf/internal/models"

type Repository interface {
	GetRating() ([]models.Rating, error)
}
