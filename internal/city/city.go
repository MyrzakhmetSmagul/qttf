package city

import (
	"fmt"
	"qttf/internal/models"
)

type Repository interface {
	Create(city *models.City) error
	GetById(id int) (models.City, error)
	GetIdByLink(link string) (int, error)
}

type CityService struct {
	Repository Repository
}

func NewCityService(repository Repository) *CityService {
	return &CityService{
		Repository: repository,
	}
}

func (s *CityService) ActualizeCities(cities []models.City) error {
	for i := range cities {
		err := s.Repository.Create(&cities[i])
		if err != nil {
			return fmt.Errorf("cityInsertion: %w\nCity: %+v", err, cities[i])
		}
	}

	return nil
}
