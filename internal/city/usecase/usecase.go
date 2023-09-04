package usecase

import (
	"context"
	"qttf/internal/city"
)

type cityUC struct {
	cityRepo city.Repository
}

func NewCityUseCase(cityRepo city.Repository) city.UseCase {
	return cityUC{cityRepo: cityRepo}
}

func (c cityUC) GetCities(ctx context.Context) {

}
