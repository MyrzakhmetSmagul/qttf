package usecase

import (
	"qttf/internal/city"
	"qttf/internal/models"
)

type cityUC struct {
	cityRepo city.Repository
}

func NewCityUseCase(cityRepo city.Repository) city.UseCase {
	return &cityUC{cityRepo: cityRepo}
}

func (c *cityUC) GetCities() ([]models.City, error) {
	return c.cityRepo.GetCities()
}
