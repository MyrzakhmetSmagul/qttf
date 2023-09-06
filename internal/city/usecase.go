package city

import "qttf/internal/models"

type UseCase interface {
	GetCities() ([]models.City, error)
}
