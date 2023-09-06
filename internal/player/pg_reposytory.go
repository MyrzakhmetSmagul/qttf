package player

import "qttf/internal/models"

type Repository interface {
	Create(player *models.Player) error
	GetPlayers() ([]models.Player, error)
	GetById(id int) (models.Player, error)
}
