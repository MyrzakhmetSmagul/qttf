package city

import (
	"qttf/internal/models"
)

type Repository interface {
	Create(city *models.City) error
	GetCities() ([]models.City, error)
	GetById(id int) (models.City, error)
	GetIdByLink(link string) (int, error)
}
