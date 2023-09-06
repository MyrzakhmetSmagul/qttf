package usecase

import (
	"qttf/internal/city"
	"qttf/internal/models"
	"qttf/internal/player"
)

type playerUC struct {
	playerRepo player.Repository
	cityRepo   city.Repository
}

func NewPlayerUseCase(playerRepo player.Repository) player.UseCase {
	return &playerUC{playerRepo: playerRepo}
}

// GetById implements player.UseCase.
func (p *playerUC) GetById(id int) (models.Player, error) {
	player, err := p.playerRepo.GetById(id)
	if err != nil {
		return player, err
	}

	player.City, err = p.cityRepo.GetById(player.City.Id)
	if err != nil {
		return player, err
	}

	return player, nil
}

// GetPlayers implements player.UseCase.
func (p *playerUC) GetPlayers() ([]models.Player, error) {
	players, err := p.playerRepo.GetPlayers()
	if err != nil {
		return nil, err
	}

	cities, err := p.cityRepo.GetCities()
	if err != nil {
		return nil, err
	}

	citiesMap := makeMap(cities)
	for i := 0; i < len(players); i++ {
		players[i].City = citiesMap[players[i].City.Id]
	}

	return players, nil
}

func makeMap(v []models.City) (cities map[int]models.City) {
	for _, v := range v {
		cities[v.Id] = v
	}
	return cities
}
