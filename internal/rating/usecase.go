package rating

import "qttf/internal/models"

type UseCase interface {
	GetRating() ([]models.Rating, error)
}
