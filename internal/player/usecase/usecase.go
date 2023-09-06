package usecase

import (
	"qttf/internal/models"
	"qttf/internal/player"
)

type playerUC struct {
	playerRepo player.Repository
}

func NewPlayerUseCase(playerRepo player.Repository) player.UseCase {
	return &playerUC{playerRepo: playerRepo}
}

// GetById implements player.UseCase.
func (p *playerUC) GetById(id int) (models.Player, error) {
	return p.playerRepo.GetById(id)
}

// GetPlayers implements player.UseCase.
func (p *playerUC) GetPlayers() ([]models.Player, error) {
	return p.playerRepo.GetPlayers()
}
