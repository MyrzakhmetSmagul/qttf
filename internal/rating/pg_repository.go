package rating

import "qttf/internal/models"

type Repository interface {
	Create(rating *models.Rating) error
	GetRating() ([]models.Rating, error)
}
